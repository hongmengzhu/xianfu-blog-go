package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamResourceGroupRelationRepository))

	gs.Provide(new(support.BaseService[RamResourceGroupRelationRepository]))
}

type RamResourceGroupRelationRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceGroupRelationEntity, int64]
	//
}

func (c *RamResourceGroupRelationRepository) FindByMark(ctx context.Context, code string) (info *entityRam.RamResourceGroupRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", code).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamResourceGroupRelationRepository) DeleteByTypeCategoryAndTypeValue(ctx context.Context, typeCategory, typeValue string) error {
	tx := c.DbModel().WithContext(ctx).Where("type_category = ?", typeCategory).Where("type_value = ?", typeValue).Delete(&entityRam.RamResourceGroupRelationEntity{})
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return tx.Error
	}
	return nil
}
