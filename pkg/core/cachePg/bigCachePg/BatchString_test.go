package bigCachePg

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"go-spring.org/spring/gs"
	StarterBigCache "go-spring.org/starter-bigcache"
)

// newTestBatchString 直接构造一个 bigcache 实例，不经过容器：
// 这些用例要验证的是 BatchString 自己的语义（TTL 信封、Hash/Set 模拟），不是 DI。
func newTestBatchString(t *testing.T) *BatchString {
	t.Helper()
	conf := bigcache.DefaultConfig(time.Hour)
	conf.Shards = 16
	conf.Verbose = false
	conf.MaxEntriesInWindow = 1024
	conf.MaxEntrySize = 1024
	conf.CleanWindow = time.Second

	client, err := bigcache.New(context.Background(), conf)
	if err != nil {
		t.Fatalf("bigcache.New: %v", err)
	}
	cache, err := StarterBigCache.NewCache(client, "rdsPg-test")
	if err != nil {
		_ = client.Close()
		t.Fatalf("NewCache: %v", err)
	}
	t.Cleanup(func() { _ = cache.Close() })
	return NewBatchString(cache)
}

func TestSetGet(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)

	if err := bs.Set(ctx, "str", "value", 0); err != nil {
		t.Fatalf("Set string: %v", err)
	}
	if got, ok := bs.Get(ctx, "str"); !ok || got != "value" {
		t.Fatalf("Get string = %q, %v; want \"value\", true", got, ok)
	}

	// []byte 与数字按原来 go-redis 的取值规则处理。
	if err := bs.Set(ctx, "bytes", []byte("raw"), 0); err != nil {
		t.Fatalf("Set bytes: %v", err)
	}
	if got, ok := bs.Get(ctx, "bytes"); !ok || got != "raw" {
		t.Fatalf("Get bytes = %q, %v; want \"raw\", true", got, ok)
	}
	if err := bs.Set(ctx, "int", 42, 0); err != nil {
		t.Fatalf("Set int: %v", err)
	}
	if got, ok := bs.Get(ctx, "int"); !ok || got != "42" {
		t.Fatalf("Get int = %q, %v; want \"42\", true", got, ok)
	}

	if _, ok := bs.Get(ctx, "missing"); ok {
		t.Fatal("Get missing = ok; want a miss")
	}
}

// TestPerEntryTTL 验证「逐条过期」：bigcache 本身只有实例级 life-window，
// 过期判定由值信封承担。
func TestPerEntryTTL(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)

	if err := bs.Set(ctx, "short", "v", 20*time.Millisecond); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, ok := bs.Get(ctx, "short"); !ok {
		t.Fatal("entry must be readable before its deadline")
	}
	time.Sleep(40 * time.Millisecond)
	if _, ok := bs.Get(ctx, "short"); ok {
		t.Fatal("entry must be a miss after its deadline")
	}
	// 过期条目在读路径上被顺手删除：clean-window=0 时没有后台 cleaner，
	// 不做这一步它会一直占着内存到容量淘汰。
	if n := bs.cache.Len(); n != 0 {
		t.Fatalf("cache still holds %d entries after the expired read; want 0", n)
	}

	// expire<=0 表示永不过期。
	if err := bs.Set(ctx, "forever", "v", 0); err != nil {
		t.Fatalf("Set: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	if got, ok := bs.Get(ctx, "forever"); !ok || got != "v" {
		t.Fatalf("Get forever = %q, %v; want \"v\", true", got, ok)
	}
}

func TestSetPipelineAndGetAllByKeys(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)

	bs.SetPipeline(ctx, map[string]any{
		"a": "1",
		"b": []byte("2"),
	})
	bs.SetPipelineTimeDuration(ctx, map[string]any{"c": "3"}, 20*time.Millisecond)

	values, ok := bs.GetAllByKeys(ctx, []string{"a", "b", "missing", "c"})
	if !ok {
		t.Fatal("GetAllByKeys = false; want true")
	}
	want := []any{"1", "2", nil, "3"}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("GetAllByKeys = %v; want %v", values, want)
	}

	// 带 TTL 的那条在过期后应变成 nil。
	time.Sleep(40 * time.Millisecond)
	values, _ = bs.GetAllByKeys(ctx, []string{"a", "c"})
	if !reflect.DeepEqual(values, []any{"1", nil}) {
		t.Fatalf("GetAllByKeys after expiry = %v; want [1 <nil>]", values)
	}
}

func TestHash(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)

	bs.HSetPipeline(ctx, "hash:1", map[string]any{"f1": "v1", "f2": []byte("v2")})
	bs.HSetPipelineMapAll(ctx, map[string]map[string]any{
		"hash:1": {"f3": "v3"},
		"hash:2": {"g1": "w1"},
	})

	all, ok := bs.HGetAll(ctx, "hash:1")
	if !ok {
		t.Fatal("HGetAll hash:1 = false; want true")
	}
	want := map[string]string{"f1": "v1", "f2": "v2", "f3": "v3"}
	if !reflect.DeepEqual(all, want) {
		t.Fatalf("HGetAll = %v; want %v", all, want)
	}

	if _, ok := bs.HGetAll(ctx, "hash:missing"); ok {
		t.Fatal("HGetAll on a missing hash = true; want false")
	}

	pipeline, ok := bs.HGetAllPipeline(ctx, []string{"hash:1", "hash:2", "hash:missing", " "})
	if !ok {
		t.Fatal("HGetAllPipeline = false; want true")
	}
	if !reflect.DeepEqual(pipeline["hash:1"], want) {
		t.Fatalf("HGetAllPipeline hash:1 = %v; want %v", pipeline["hash:1"], want)
	}
	if !reflect.DeepEqual(pipeline["hash:2"], map[string]string{"g1": "w1"}) {
		t.Fatalf("HGetAllPipeline hash:2 = %v", pipeline["hash:2"])
	}
	// 不存在的 Hash 与 Redis 的 HGETALL 一样返回空 map。
	if got, exists := pipeline["hash:missing"]; !exists || len(got) != 0 {
		t.Fatalf("HGetAllPipeline missing = %v, %v; want empty map, true", got, exists)
	}
}

func TestSetAndGetAllBySetKey(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)

	// 重复成员只保留一份。
	if err := bs.SAdd(ctx, "set:1", []string{"a", "b", "a"}); err != nil {
		t.Fatalf("SAdd: %v", err)
	}
	members, ok := bs.SMembers(ctx, "set:1")
	if !ok || !reflect.DeepEqual(members, []string{"a", "b"}) {
		t.Fatalf("SMembers = %v, %v; want [a b], true", members, ok)
	}
	if _, ok := bs.SMembers(ctx, "set:missing"); ok {
		t.Fatal("SMembers on a missing set = true; want false")
	}

	bs.SetPipeline(ctx, map[string]any{"a": "va", "b": "vb"})
	// c 没有对应的值：GetAllBySetKey 应跳过它，而不是返回空值。
	if err := bs.SAdd(ctx, "set:1", []string{"c"}); err != nil {
		t.Fatalf("SAdd: %v", err)
	}

	pairs, ok := bs.GetAllBySetKey(ctx, "set:1")
	if !ok {
		t.Fatal("GetAllBySetKey = false; want true")
	}
	want := []any{"a", "va", "b", "vb"}
	if !reflect.DeepEqual(pairs, want) {
		t.Fatalf("GetAllBySetKey = %v; want %v", pairs, want)
	}

	if _, ok := bs.GetAllBySetKey(ctx, "set:missing"); ok {
		t.Fatal("GetAllBySetKey on a missing set = true; want false (same as the old Lua failure)")
	}
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)

	if err := bs.Set(ctx, "k1", "v", 0); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := bs.Set(ctx, "k2", "v", 0); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := bs.Del(ctx, "k1", "k2", "never-existed"); err != nil {
		t.Fatalf("Del: %v", err)
	}
	if _, ok := bs.Get(ctx, "k1"); ok {
		t.Fatal("k1 must be gone after Del")
	}
}

// TestWiring 用 go-spring 的容器测试钩子验证真实装配：
// starter-bigcache 按 spring.bigcache.instances.app 建实例，driver=longterm 选中
// bigCachePg 注册的 Driver，BatchString 按名字拿到同一个实例。
func TestWiring(t *testing.T) {
	gs.Web(false).Configure(func(app gs.App) {
		// longterm Driver 要求 life-window>=24h、clean-window=0、有内存上限。
		app.Property("spring.bigcache.default.shards", "16")
		app.Property("spring.bigcache.default.life-window", "720h")
		app.Property("spring.bigcache.default.clean-window", "0s")
		app.Property("spring.bigcache.default.max-entries-in-window", "1024")
		app.Property("spring.bigcache.default.max-entry-size", "1024")
		app.Property("spring.bigcache.default.hard-max-cache-size", "16")
		app.Property("spring.bigcache.instances.app.life-window", "720h")
		app.Property("spring.bigcache.instances.app.driver", "longterm")
	}).RunTest(t, func(ts *struct {
		Batch *BatchString           `autowire:"?"`
		Cache *StarterBigCache.Cache `autowire:"app"`
	}) {
		if ts.Cache == nil {
			t.Fatal("expected the bigcache instance bean named \"app\"")
		}
		if ts.Batch == nil {
			t.Fatal("expected the BatchString bean")
		}
		if ts.Batch.cache != ts.Cache {
			t.Fatal("BatchString must inject the same \"app\" instance, not a second cache")
		}

		ctx := context.Background()
		if err := ts.Batch.Set(ctx, "wired", "ok", 0); err != nil {
			t.Fatalf("Set through the container-wired BatchString: %v", err)
		}
		if got, ok := ts.Batch.Get(ctx, "wired"); !ok || got != "ok" {
			t.Fatalf("Get = %q, %v; want \"ok\", true", got, ok)
		}
	})
}
