package configModel

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"go-spring.org/log"
)

// DeleteModelData 删除模型关联数据：模型字段、字段规则、模型级规则、配置行、配置列表行，并清理配置缓存。
func (s *Sp) DeleteModelData(ctx context.Context, info *entityBasic.BasicConfigModelEntity) {
	if s == nil || info == nil {
		return
	}
	if err := s.repConfig.DeleteAllByModelNo(ctx, info.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete model config error")
	}
	if err := s.repConfigList.DeleteAllByModelNo(ctx, info.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete model config list error")
	}
	if err := s.repField.DeleteAllByModelNo(ctx, info.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete model fields error")
	}
	if err := s.repFieldRules.DeleteAllByValueNo(ctx, info.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete field rules error")
	}
	if err := s.repRule.DeleteAllByModelNo(ctx, info.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete model rules error")
	}
	if s.cache != nil {
		if err := s.cache.Refresh(ctx, info.TenantNo, info.Model); err != nil {
			log.Error(ctx, log.TagAppDef, err, "refresh config cache error")
		}
	}
}
