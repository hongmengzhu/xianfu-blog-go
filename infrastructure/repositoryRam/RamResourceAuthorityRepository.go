package repositoryRam

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(RamResourceAuthorityRepository))

	gs.Provide(new(support.BaseService[RamResourceAuthorityRepository]))
}

type RamResourceAuthorityRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceAuthorityEntity, int64]
	//
}

func (c *RamResourceAuthorityRepository) FindByMark(ctx context.Context, code string) (info *entityRam.RamResourceAuthorityEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", code).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamResourceAuthorityRepository) FindAllByGroupIdStringIn(ctx context.Context, code []string) (info []*entityRam.RamResourceAuthorityEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("group_id in ?", code).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *RamResourceAuthorityRepository) FindAllByTypeCategoryAndGroupIdStringIn(ctx context.Context, typeCategory string, code []string) (info []*entityRam.RamResourceAuthorityEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_category = ?", typeCategory).Where("group_id in ?", code).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamResourceAuthorityRepository) DeleteByMark(ctx context.Context, code string) error {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", code).Delete(&entityRam.RamResourceAuthorityEntity{})
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return tx.Error
	}
	return nil
}
