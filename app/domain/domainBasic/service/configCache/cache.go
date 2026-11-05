package configCache

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/cache/basic/modCacheBasicConfig"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/basic/key/basicConfigKey"
	"go-spring.org/log"
)

// CacheTTL 配置缓存兜底过期时间；正常写操作会立即 Refresh，TTL 只用于异常场景自愈。
const CacheTTL = 30 * time.Minute

// Key 返回配置缓存键：tenant_no + model。
func (c *Sp) Key(tenantNo, model string) string {
	return basicConfigKey.ConfigByTenantNoAndModel(tenantNo, model)
}

// Get 读取缓存；未命中时 ok=false，由调用方回源数据库并调用 Refresh。
func (c *Sp) Get(ctx context.Context, tenantNo, model string) (items []modCacheBasicConfig.Item, ok bool) {
	if c == nil || c.Rdt == nil {
		return nil, false
	}
	key := c.Key(tenantNo, model)
	raw, exists := c.Rdt.Get(ctx, key)
	if !exists {
		return nil, false
	}
	items = make([]modCacheBasicConfig.Item, 0)
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		log.Error(ctx, log.TagAppDef, err, "配置缓存反序列化失败", log.String("key", key))
		_ = c.Rdt.Del(ctx, key)
		return nil, false
	}
	return items, true
}

// Refresh 按 tenant_no + model 重建缓存；无启用配置时删除 key。
func (c *Sp) Refresh(ctx context.Context, tenantNo, model string) error {
	if c == nil || c.Rdt == nil || tenantNo == "" || model == "" {
		return nil
	}
	infos, _ := c.RepConfig.FindAllEnabledByTenantNoAndModel(ctx, tenantNo, model)
	items := make([]modCacheBasicConfig.Item, 0, len(infos))
	for _, info := range infos {
		items = append(items, CopyItem(info))
	}
	return c.Set(ctx, tenantNo, model, items)
}

// Set 直接把已构建好的配置项写入缓存；空数组也写入，作为负缓存避免反复回源。
func (c *Sp) Set(ctx context.Context, tenantNo, model string, items []modCacheBasicConfig.Item) error {
	if c == nil || c.Rdt == nil || tenantNo == "" || model == "" {
		return nil
	}
	key := c.Key(tenantNo, model)
	if items == nil {
		items = []modCacheBasicConfig.Item{}
	}
	payload, err := json.Marshal(items)
	if err != nil {
		return err
	}
	// BatchString.Set 的 expire 为 time.Duration；传 0 表示不设置逐条过期，由实例级 TTL 兜底。
	return c.Rdt.Set(ctx, key, string(payload), 0)
}

// Override 全量重建所有租户/模型的配置缓存。
func (c *Sp) Override(ctx context.Context) error {
	if c == nil || c.Rdt == nil {
		return nil
	}
	infos, _ := c.RepConfig.FindAllEnabled(ctx)
	group := make(map[string][]modCacheBasicConfig.Item)
	for _, info := range infos {
		if info.TenantNo == "" || info.Model == "" {
			continue
		}
		group[info.TenantNo+"\x00"+info.Model] = append(group[info.TenantNo+"\x00"+info.Model], CopyItem(info))
	}
	for groupKey, items := range group {
		var tenantNo, model string
		for i, ch := range groupKey {
			if ch == '\x00' {
				tenantNo = groupKey[:i]
				model = groupKey[i+1:]
				break
			}
		}
		payload, err := json.Marshal(items)
		if err != nil {
			return err
		}
		if err := c.Rdt.Set(ctx, c.Key(tenantNo, model), string(payload), 0); err != nil {
			return err
		}
	}
	return nil
}

// CopyItem 把配置实体转换为缓存项。
func CopyItem(info *entityB72.B72ConfigEntity) modCacheBasicConfig.Item {
	return modCacheBasicConfig.Item{
		ID:           info.ID,
		No:           info.No,
		TenantNo:     info.TenantNo,
		ModelNo:      info.ModelNo,
		Model:        info.Model,
		Module:       info.Module,
		ModuleSub:    info.ModuleSub,
		Field:        info.Field,
		FieldPath:    info.FieldPath,
		Name:         info.Name,
		Description:  info.Description,
		Show:         info.Show,
		Binary:       info.Binary,
		ValueType:    info.ValueType,
		FormCode:     info.FormCode,
		DefaultValue: info.DefaultValue,
		Value:        info.Value,
		ValueBinary:  info.ValueBinary,
		CategoryNo:   info.CategoryNo,
		Sort:         info.Sort,
	}
}
