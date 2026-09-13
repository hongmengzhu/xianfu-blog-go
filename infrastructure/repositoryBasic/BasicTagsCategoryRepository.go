package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicTagsCategoryRepository))

	gs.Provide(new(support.BaseService[BasicTagsCategoryRepository]))
}

type BasicTagsCategoryRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicTagsCategoryEntity, int64]
}

func (c *BasicTagsCategoryRepository) FindAllByNoLinkAndTypeSys(ctx context.Context, code string, tpSys string) (info []*entityBasic.BasicTagsCategoryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_sys = ?", tpSys).Where("no_link like ?", "%"+code+"%").Find(&info)
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
