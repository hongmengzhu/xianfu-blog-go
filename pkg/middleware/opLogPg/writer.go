package opLogPg

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

const (
	channelSize     = 4096
	batchSize       = 100
	flushInterval   = time.Second
	writeTimeout    = 3 * time.Second
	retryTimes      = 3
	shutdownTimeout = 5 * time.Second
	fallbackDir     = "logs/oplog-fallback"
)

// defaultWriter 进程级写入器，供 Global 中间件在请求期使用。
var defaultWriter atomic.Pointer[Writer]

func init() {
	gs.Provide(NewWriter).Export(gs.As[gs.Server]())
}

// Writer 操作日志异步批量写入器。
type Writer struct {
	repo *repositoryBasic.BasicOperationLogRepository `autowire:"?"`

	ch chan *entityBasic.BasicOperationLogEntity

	mu      sync.RWMutex
	closed  bool
	wg      sync.WaitGroup
	stopOne sync.Once
	fbMu    sync.Mutex
}

// NewWriter 构造写入器并登记为进程级默认实例。
func NewWriter(repo *repositoryBasic.BasicOperationLogRepository) *Writer {
	w := &Writer{
		repo: repo,
		ch:   make(chan *entityBasic.BasicOperationLogEntity, channelSize),
	}
	defaultWriter.Store(w)
	return w
}

// Run 实现 gs.Server：启动 worker，停机时排空队列。
func (w *Writer) Run(ctx context.Context, sig gs.ReadySignal) error {
	w.wg.Add(1)
	go w.loop(ctx)

	<-sig.TriggerAndWait()
	<-ctx.Done()
	w.stop()
	// 若 Stop 在 Run 之前被调用，stop() 无法等待后来启动的 worker，这里再等一次。
	w.wg.Wait()
	return nil
}

// Stop 实现 gs.Server：幂等关闭。
func (w *Writer) Stop(ctx context.Context) error {
	w.stop()
	return nil
}

func (w *Writer) stop() {
	w.stopOne.Do(func() {
		w.closeChannel()
		done := make(chan struct{})
		go func() {
			w.wg.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(shutdownTimeout):
			log.Error(context.Background(), log.TagAppDef, errors.New("操作日志队列排空超时"), "oplog writer shutdown timeout")
		}
	})
}

// closeChannel 停止接收新日志，之后入队走同步兜底。
func (w *Writer) closeChannel() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return
	}
	w.closed = true
	close(w.ch)
}

// enqueue 非阻塞入队；队列满或已关闭时落兜底文件。
func (w *Writer) enqueue(entity *entityBasic.BasicOperationLogEntity) {
	w.mu.RLock()
	if w.closed {
		w.mu.RUnlock()
		w.writeOneSync(context.Background(), entity)
		return
	}
	select {
	case w.ch <- entity:
		w.mu.RUnlock()
	default:
		w.mu.RUnlock()
		w.fallback(entity)
	}
}

func (w *Writer) loop(ctx context.Context) {
	defer w.wg.Done()
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	buffer := make([]*entityBasic.BasicOperationLogEntity, 0, batchSize)
	flush := func() {
		if len(buffer) == 0 {
			return
		}
		logs := buffer
		buffer = make([]*entityBasic.BasicOperationLogEntity, 0, batchSize)
		w.flush(logs)
	}

	for {
		select {
		case entity, ok := <-w.ch:
			if !ok {
				flush()
				return
			}
			buffer = append(buffer, entity)
			if len(buffer) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			// channel 由 stop() 关闭，继续读到关闭为止。
			for entity := range w.ch {
				buffer = append(buffer, entity)
				if len(buffer) >= batchSize {
					flush()
				}
			}
			flush()
			return
		}
	}
}

func (w *Writer) flush(logs []*entityBasic.BasicOperationLogEntity) {
	// 用 Background，避免停机时 ctx 已取消导致落库失败而全部转兜底文件。
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	if err := w.persist(ctx, logs); err != nil {
		log.Error(ctx, log.TagAppDef, err, "操作日志批量写入失败，转兜底")
	}
}

// persist 批量写；失败按退避重试，仍失败逐条写，坏行落兜底。
func (w *Writer) persist(ctx context.Context, logs []*entityBasic.BasicOperationLogEntity) error {
	if len(logs) == 0 {
		return nil
	}
	if w.repo == nil {
		for _, item := range logs {
			w.fallback(item)
		}
		return errors.New("oplog repository 未注入")
	}

	var lastErr error
	for attempt := range retryTimes {
		if attempt > 0 {
			select {
			case <-ctx.Done():
			case <-time.After(time.Duration(attempt) * 300 * time.Millisecond):
			}
		}
		if err := w.repo.CreateLogsBatch(ctx, logs, batchSize); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	// 逐条重试，定位坏行。
	var failed int
	for _, item := range logs {
		if err := w.repo.CreateLogsBatch(ctx, []*entityBasic.BasicOperationLogEntity{item}, 1); err != nil {
			failed++
			w.fallback(item)
		}
	}
	if failed > 0 {
		return lastErr
	}
	return nil
}

func (w *Writer) writeOneSync(ctx context.Context, entity *entityBasic.BasicOperationLogEntity) {
	if w.repo == nil {
		w.fallback(entity)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	if err := w.repo.CreateLogsBatch(ctx, []*entityBasic.BasicOperationLogEntity{entity}, 1); err != nil {
		w.fallback(entity)
	}
}

// fallback 写本地 JSONL 兜底文件，供后续补偿。
func (w *Writer) fallback(entity *entityBasic.BasicOperationLogEntity) {
	if entity == nil {
		return
	}
	bs, err := json.Marshal(entity)
	if err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "操作日志兜底序列化失败")
		return
	}
	w.fbMu.Lock()
	defer w.fbMu.Unlock()
	if err := os.MkdirAll(fallbackDir, 0o755); err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "创建操作日志兜底目录失败")
		return
	}
	name := filepath.Join(fallbackDir, "oplog-fallback-"+time.Now().Format("20060102")+".jsonl")
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "打开操作日志兜底文件失败")
		return
	}
	defer f.Close()
	if _, err := f.Write(append(bs, '\n')); err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "写操作日志兜底文件失败")
	}
}
