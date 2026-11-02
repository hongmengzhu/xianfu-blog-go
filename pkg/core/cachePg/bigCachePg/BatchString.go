package bigCachePg

// BatchString 是原来的 Redis 批量缓存门面，现在整体改为进程内缓存（starter-bigcache）。
//
// 为什么不再用 Redis：
//   - 本项目的缓存都是「可重建」的读缓存，数据源始终是数据库，进程内缓存省掉一次网络往返；
//   - bigcache 没有 pipeline / Hash / Set / Lua，所以下面每个方法都换成 bigcache 能做到的
//     等价实现，语义按调用方的真实用法收敛：
//     · SetPipeline / SetPipelineTimeDuration：逐个 Set。bigcache 是进程内操作，没有网络
//       请求，本来就不需要 pipeline，循环写入就是最省的形态；
//     · 逐条过期：bigcache 只有「实例级 life-window」，没有 per-key TTL。这里在值前面加一个
//       12 字节信封（magic + 绝对过期时间），Get 时判断是否过期，因此每个 key 仍然可以独立过期；
//     · Hash：bigcache 只有 string->[]byte，Hash 用一个 key 存 JSON 对象来模拟，
//       读改写用分段锁保护；
//     · Set / Lua(SMEMBERS+MGET)：用同样的 JSON 数组模拟集合，GetAllBySetKey 复刻
//       原来 SMEMBERS+MGET 脚本的返回形状 [key1, value1, key2, value2, ...]。
//
// 缓存实例来自 starter-bigcache：bean 名就是配置里的实例名（spring.bigcache.instances.<name>）。
// 这里固定注入 "app" 实例，配置见 data/config/app.toml。

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/fnv"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	StarterBigCache "go-spring.org/starter-bigcache"
	"go-spring.org/stdlib/errutil"
)

// BigCacheInstance 是本包使用的 starter-bigcache 实例名，必须与
// spring.bigcache.instances.<name> 保持一致；bean 名同样是这个名字。
const BigCacheInstance = "app"

func init() {
	gs.Provide(new(BatchString))
}

type BatchString struct {
	// 按名字注入：bigcache 的 bean 名就是实例名，按类型注入会因多实例而歧义。
	cache *StarterBigCache.Cache `autowire:"app"`

	// Hash / Set 是「读-改-写」，需要保证同一个 key 的并发更新不丢数据。
	// 分段锁按 key 的哈希取模，避免所有 key 抢同一把锁。
	locks [lockStripes]sync.Mutex
}

// lockStripes 是分段锁的数量，必须是 2 的幂（这里用位与取模）。
const lockStripes = 64

// NewBatchString 供手工构造使用；容器里走 gs.Provide(new(BatchString)) 的字段注入。
func NewBatchString(cache *StarterBigCache.Cache) *BatchString {
	return &BatchString{cache: cache}
}

// lockFor 返回 key 对应的分段锁。
func (t *BatchString) lockFor(key string) *sync.Mutex {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return &t.locks[h.Sum32()&(lockStripes-1)]
}

// ----------------------------------------------------------------------------
// 值信封：bigcache 没有 per-key TTL，用 12 字节头把过期时间写进值里
// ----------------------------------------------------------------------------

const (
	// 4 字节 magic + 8 字节绝对过期时间（UnixNano，0 表示永不过期）
	valueHeaderSize = 4 + 8
)

var valueMagic = [4]byte{'B', 'C', 'V', '1'}

// encodeValue 给 payload 打上过期信封。expire <= 0 表示永不过期（与 Redis 的 0 语义一致）。
// 过期时间在写入时算成绝对时间，所以同一批写入共用同一个 now。
func encodeValue(payload []byte, expire time.Duration, now time.Time) []byte {
	var deadline int64
	if expire > 0 {
		deadline = now.Add(expire).UnixNano()
	}
	buf := make([]byte, valueHeaderSize+len(payload))
	copy(buf[:4], valueMagic[:])
	binary.BigEndian.PutUint64(buf[4:valueHeaderSize], uint64(deadline))
	copy(buf[valueHeaderSize:], payload)
	return buf
}

// decodeValue 拆掉信封；ok=false 表示已过期，调用方按未命中处理。
// 没有 magic 的值原样返回，兼容直接通过 Cache.Client 写入的裸值。
func decodeValue(entry []byte, now time.Time) (payload []byte, ok bool) {
	if len(entry) < valueHeaderSize || !bytes.Equal(entry[:4], valueMagic[:]) {
		return entry, true
	}
	deadline := int64(binary.BigEndian.Uint64(entry[4:valueHeaderSize]))
	if deadline != 0 && !now.Before(time.Unix(0, deadline)) {
		return nil, false
	}
	return entry[valueHeaderSize:], true
}

// encodeAny 把调用方传入的 any 变成 []byte，取值规则尽量贴近原来的 go-redis：
// string / []byte 原样，数字与 bool 转成文本，其余（结构体、map 等）走 JSON。
func encodeAny(value any) ([]byte, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		return []byte(rv.String()), nil
	case reflect.Bool:
		return strconv.AppendBool(nil, rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(nil, rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(nil, rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.AppendFloat(nil, rv.Float(), 'g', -1, 32), nil
	case reflect.Float64:
		return strconv.AppendFloat(nil, rv.Float(), 'g', -1, 64), nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, errutil.Explain(err, "bigcache: marshal value of type %T", value)
	}
	return data, nil
}

// getPayload 读取 key 并拆掉信封；ok=false 表示不存在或已过期。
//
// 过期时顺手删除：clean-window=0 时 bigcache 的后台 cleaner 不会跑，读路径自己回收，
// 过期条目才不会一直占着内存等到容量淘汰。
func (t *BatchString) getPayload(ctx context.Context, key string) ([]byte, bool) {
	raw, err := t.cache.Get(ctx, key)
	if err != nil {
		if !errors.Is(err, bigcache.ErrEntryNotFound) {
			log.Error(ctx, log.TagAppDef, err, "获取缓存失败: key="+key)
		}
		return nil, false
	}
	payload, ok := decodeValue(raw, time.Now())
	if !ok {
		if err := t.cache.Delete(ctx, key); err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
			log.Warn(ctx, log.TagAppDef, "删除过期缓存失败: key="+key, log.Any("err", err.Error()))
		}
		return nil, false
	}
	return payload, true
}

// ----------------------------------------------------------------------------
// 单 key 读写
// ----------------------------------------------------------------------------

// Set 写入一个 key，expire<=0 表示不设置过期时间。
// 对应原来调用方直接使用的 rdb.Set(ctx, key, value, ttl)。
func (t *BatchString) Set(ctx context.Context, key string, value any, expire time.Duration) error {
	payload, err := encodeAny(value)
	if err != nil {
		return err
	}
	if err := t.cache.Set(ctx, key, encodeValue(payload, expire*time.Second, time.Now())); err != nil {
		return errutil.Explain(err, "bigcache: set key %q", key)
	}
	return nil
}

// Get 读取一个 key；返回 (值, 是否存在)。已过期与不存在都返回 false。
func (t *BatchString) Get(ctx context.Context, key string) (string, bool) {
	payload, ok := t.getPayload(ctx, key)
	if !ok {
		return "", false
	}
	return string(payload), true
}

// Del 删除一批 key。删除不存在的 key 不算错误。
func (t *BatchString) Del(ctx context.Context, keys ...string) error {
	var firstErr error
	for _, key := range keys {
		if err := t.cache.Delete(ctx, key); err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
			if firstErr == nil {
				firstErr = errutil.Explain(err, "bigcache: delete key %q", key)
			}
		}
	}
	return firstErr
}

// ----------------------------------------------------------------------------
// 批量读写
// ----------------------------------------------------------------------------

// SetPipeline 批量写入，不设置过期时间。
//
// bigcache 没有 pipeline，也不存在网络往返：逐个 Set 已经是本地操作里最省的形态。
func (t *BatchString) SetPipeline(ctx context.Context, keysValues map[string]any) {
	t.setMany(ctx, keysValues, 0)
}

// SetPipelineTimeDuration 批量写入，并给每个 key 设置相同的过期时间。
func (t *BatchString) SetPipelineTimeDuration(ctx context.Context, keysValues map[string]any, expire time.Duration) {
	t.setMany(ctx, keysValues, expire)
}

func (t *BatchString) setMany(ctx context.Context, keysValues map[string]any, expire time.Duration) {
	if len(keysValues) == 0 {
		return
	}
	// 同一批写入共用一个 now，避免每条都取一次时钟。
	now := time.Now()
	failed := 0
	for key, value := range keysValues {
		payload, err := encodeAny(value)
		if err != nil {
			failed++
			log.Error(ctx, log.TagAppDef, err, "批量操作失败: key="+key)
			continue
		}
		if err := t.cache.Set(ctx, key, encodeValue(payload, expire, now)); err != nil {
			failed++
			log.Error(ctx, log.TagAppDef, err, "批量操作失败: key="+key)
		}
	}
	log.Debug(ctx, log.TagAppDef, "批量操作命令数:", func() []log.Field {
		return []log.Field{
			log.Int("keys", len(keysValues)),
			log.Int("failed", failed),
			log.String("expire", expire.String()),
		}
	})
}

// GetAllByKeys 批量读取，返回值与 keys 一一对应；不存在的 key 对应 nil。
// 对应原 Redis 的 MGet / BATCH_MGET 脚本。
func (t *BatchString) GetAllByKeys(ctx context.Context, keys []string) ([]any, bool) {
	result := make([]any, len(keys))
	for i, key := range keys {
		if payload, ok := t.getPayload(ctx, key); ok {
			result[i] = string(payload)
		}
	}
	return result, true
}

// ----------------------------------------------------------------------------
// 集合：模拟 Redis 的 SAdd / SMembers / SMEMBERS+MGET 脚本
// ----------------------------------------------------------------------------

// SAdd 向集合追加成员，返回 nil 表示成功。集合以 JSON 数组存在同一个 key 下。
func (t *BatchString) SAdd(ctx context.Context, key string, members []string) error {
	members = compactStrings(members)
	if len(members) == 0 {
		return nil
	}
	lock := t.lockFor(key)
	lock.Lock()
	defer lock.Unlock()

	current, _ := t.loadStrings(ctx, key)
	seen := make(map[string]struct{}, len(current)+len(members))
	merged := make([]string, 0, len(current)+len(members))
	for _, item := range current {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		merged = append(merged, item)
	}
	for _, item := range members {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		merged = append(merged, item)
	}
	return t.storeStrings(ctx, key, merged)
}

// SMembers 返回集合全部成员；集合不存在时返回 (nil, false)。
func (t *BatchString) SMembers(ctx context.Context, key string) ([]string, bool) {
	return t.loadStrings(ctx, key)
}

// GetAllBySetKey 复刻原来 SMEMBERS+MGET 脚本的行为：
// 读取集合 key 的所有成员，再批量取每个成员的值，返回 [key1, value1, key2, value2, ...]。
// 集合不存在或为空时返回 (nil, false)，与原来 Lua 执行失败的表现一致。
func (t *BatchString) GetAllBySetKey(ctx context.Context, setKey string) ([]any, bool) {
	members, ok := t.SMembers(ctx, setKey)
	if !ok || len(members) == 0 {
		return nil, false
	}
	result := make([]any, 0, len(members)*2)
	for _, member := range members {
		if payload, ok := t.getPayload(ctx, member); ok {
			result = append(result, member, string(payload))
		}
	}
	if len(result) == 0 {
		return nil, false
	}
	return result, true
}

// loadStrings 读取一个 JSON 字符串数组。ok=false 表示 key 不存在或已过期。
func (t *BatchString) loadStrings(ctx context.Context, key string) ([]string, bool) {
	payload, ok := t.getPayload(ctx, key)
	if !ok {
		return nil, false
	}
	var items []string
	if err := json.Unmarshal(payload, &items); err != nil {
		// 同一个 key 上写过的不是集合（例如普通字符串），按不存在处理，
		// 让 SAdd 从空集合开始，而不是把已有数据当集合拼坏。
		return nil, false
	}
	return items, true
}

// storeStrings 把字符串数组写成 JSON 并保存，集合本身不设置过期时间。
func (t *BatchString) storeStrings(ctx context.Context, key string, items []string) error {
	data, err := json.Marshal(items)
	if err != nil {
		return errutil.Explain(err, "bigcache: marshal set %q", key)
	}
	if err := t.cache.Set(ctx, key, encodeValue(data, 0, time.Now())); err != nil {
		return errutil.Explain(err, "bigcache: store set %q", key)
	}
	return nil
}

// compactStrings 去掉空字符串，并保持原顺序。
func compactStrings(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if strPg.IsBlank(item) {
			continue
		}
		out = append(out, item)
	}
	return out
}

// ----------------------------------------------------------------------------
// Hash：bigcache 没有 Hash，用一个 key 存 JSON 对象模拟
// ----------------------------------------------------------------------------

// HSetPipeline 批量设置同一个 Hash 的字段。
func (t *BatchString) HSetPipeline(ctx context.Context, hashKey string, keysValues map[string]any) {
	t.HSetPipelineMapAll(ctx, map[string]map[string]any{hashKey: keysValues})
}

// HSetPipelineMapAll 批量设置多个 Hash 的字段。
func (t *BatchString) HSetPipelineMapAll(ctx context.Context, keysValues map[string]map[string]any) {
	if len(keysValues) == 0 {
		return
	}
	failed := 0
	for hashKey, fields := range keysValues {
		if len(fields) == 0 {
			continue
		}
		if err := t.mergeHash(ctx, hashKey, fields); err != nil {
			failed++
			log.Error(ctx, log.TagAppDef, err, "批量操作失败: hashKey="+hashKey)
		}
	}
	log.Debug(ctx, log.TagAppDef, "批量操作命令数:", func() []log.Field {
		return []log.Field{
			log.Int("hashKeys", len(keysValues)),
			log.Int("failed", failed),
		}
	})
}

// mergeHash 在分段锁内完成「读 -> 合并 -> 写」。
func (t *BatchString) mergeHash(ctx context.Context, hashKey string, fields map[string]any) error {
	lock := t.lockFor(hashKey)
	lock.Lock()
	defer lock.Unlock()

	current, _ := t.loadHash(ctx, hashKey)
	if current == nil {
		current = make(map[string]string, len(fields))
	}
	for field, value := range fields {
		payload, err := encodeAny(value)
		if err != nil {
			return err
		}
		current[field] = string(payload)
	}
	data, err := json.Marshal(current)
	if err != nil {
		return errutil.Explain(err, "bigcache: marshal hash %q", hashKey)
	}
	if err := t.cache.Set(ctx, hashKey, encodeValue(data, 0, time.Now())); err != nil {
		return errutil.Explain(err, "bigcache: store hash %q", hashKey)
	}
	return nil
}

// HGetAll 获取 Hash 的所有字段；不存在返回 (nil, false)。
func (t *BatchString) HGetAll(ctx context.Context, hashKey string) (map[string]string, bool) {
	return t.loadHash(ctx, hashKey)
}

// HGetAllPipeline 批量获取多个 Hash。不存在的 Hash 返回空 map，
// 与 Redis 的 HGETALL 对空 key 的行为一致。
func (t *BatchString) HGetAllPipeline(ctx context.Context, hashKeys []string) (map[string]map[string]string, bool) {
	result := make(map[string]map[string]string, len(hashKeys))
	for _, hashKey := range hashKeys {
		if strPg.IsBlank(hashKey) {
			continue
		}
		fields, _ := t.loadHash(ctx, hashKey)
		if fields == nil {
			fields = map[string]string{}
		}
		result[hashKey] = fields
	}
	log.Debug(ctx, log.TagAppDef, "批量操作命令数:", func() []log.Field {
		return []log.Field{log.Int("hashKeys", len(hashKeys))}
	})
	return result, true
}

// loadHash 读取 JSON 对象。ok=false 表示 key 不存在、已过期或不是 Hash。
func (t *BatchString) loadHash(ctx context.Context, hashKey string) (map[string]string, bool) {
	payload, ok := t.getPayload(ctx, hashKey)
	if !ok {
		return nil, false
	}
	var fields map[string]string
	if err := json.Unmarshal(payload, &fields); err != nil {
		return nil, false
	}
	return fields, true
}
