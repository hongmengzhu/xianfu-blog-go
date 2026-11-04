package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicConfigCategoryRepository))

	gs.Provide(new(support.BaseService[BasicConfigCategoryRepository]))
}

type BasicConfigCategoryRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigCategoryEntity, int64]
}

// FindAllEnabledByTenantNo 查询租户下所有启用分类，按 sort、创建时间排序。
func (c *BasicConfigCategoryRepository) FindAllEnabledByTenantNo(ctx context.Context, tenantNo string) (info []*entityBasic.BasicConfigCategoryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).
		Where("tenant_no=?", tenantNo).
		Where("state=?", 1).
		Order("sort asc,create_at").
		Find(&info)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
