package cacheBm

import (
	"context"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/bm/key/bmKeyPg"
	StarterBigCache "go-spring.org/starter-bigcache"
)

// newTestBatchString 直接构造 bigcache 实例，不经过容器：
// 这些用例只验证缓存键的读写语义，不依赖 DI。
func newTestBatchString(t *testing.T) *bigCachePg.BatchString {
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
	return bigCachePg.NewBatchString(cache)
}

// seedCategory 按生产环境的写缓存方式写入分类：
// bm:bmCate:<tenantNo>:<no> 存 JSON，bm:bmCateNoCode:<tenantNo>:<code> 存 no。
func seedCategory(t *testing.T, bs *bigCachePg.BatchString, tenantNo string, category modBmBookmarkCategory.Cache) {
	t.Helper()
	ctx := context.Background()
	str, err := json.Marshal(category)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if err := bs.Set(ctx, bmKeyPg.BookmarkCategoryTenantNo(tenantNo, category.No), str, 0); err != nil {
		t.Fatalf("Set category: %v", err)
	}
	if err := bs.Set(ctx, bmKeyPg.BookmarkCategoryTenantNoAndNoByCode(tenantNo, category.Code), category.No, 0); err != nil {
		t.Fatalf("Set code mapping: %v", err)
	}
}

// TestGetByKeyByCode 页面 /bookmark/category/:code 传的是 code，
// GetByKey 必须先用 code 换成 no，再按 no 读取分类 JSON。
func TestGetByKeyByCode(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)
	c := &CacheBookmarkCategory{rdt: bs}

	tenantNo := "-1"
	want := modBmBookmarkCategory.Cache{
		No:   "2611021604269832474083",
		Code: "2611021604269836697181",
		Name: "go",
	}
	seedCategory(t, bs, tenantNo, want)

	got, ok := c.GetByKey(ctx, want.Code, tenantNo)
	if !ok || got == nil {
		t.Fatalf("GetByKey(code) = %v, %v; want the category", got, ok)
	}
	if got.No != want.No || got.Name != want.Name {
		t.Fatalf("GetByKey(code) = %+v; want %+v", *got, want)
	}
}

// TestGetByKeyByNo 兼容 URL 直接传 no 的情况。
func TestGetByKeyByNo(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)
	c := &CacheBookmarkCategory{rdt: bs}

	tenantNo := "-1"
	want := modBmBookmarkCategory.Cache{No: "2611021604269832474083", Code: "2611021604269836697181", Name: "go"}
	seedCategory(t, bs, tenantNo, want)

	got, ok := c.GetByKey(ctx, want.No, tenantNo)
	if !ok || got == nil || got.No != want.No {
		t.Fatalf("GetByKey(no) = %v, %v; want the category", got, ok)
	}
}

// TestGetByKeyMiss 未命中时返回 (nil,false)，且不能因 nil 值打印格式错误。
func TestGetByKeyMiss(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)
	c := &CacheBookmarkCategory{rdt: bs}

	if got, ok := c.GetByKey(ctx, "not-exists", "-1"); ok || got != nil {
		t.Fatalf("GetByKey(miss) = %v, %v; want (nil,false)", got, ok)
	}
}

// TestGetAllByKeysRetMapSkipsMiss 批量读取中缺失的 key 为 nil，必须跳过而不是当成格式错误。
func TestGetAllByKeysRetMapSkipsMiss(t *testing.T) {
	ctx := context.Background()
	bs := newTestBatchString(t)
	c := &CacheBookmarkCategory{rdt: bs}

	tenantNo := "-1"
	want := modBmBookmarkCategory.Cache{No: "2611021604269832474083", Code: "2611021604269836697181", Name: "go"}
	seedCategory(t, bs, tenantNo, want)

	got, ok := c.GetAllByKeysRetMap(ctx, []string{want.No, "not-exists"}, tenantNo)
	if !ok {
		t.Fatal("GetAllByKeysRetMap = false; want true")
	}
	if len(got) != 1 || got[want.No] == nil || got[want.No].Name != want.Name {
		t.Fatalf("GetAllByKeysRetMap = %+v; want only %q", got, want.No)
	}
}
