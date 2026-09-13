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
	gs.Provide(new(RamResourceGroupRepository))

	gs.Provide(new(support.BaseService[RamResourceGroupRepository]))
}

type RamResourceGroupRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceGroupEntity, int64]
	//
}

func (c *RamResourceGroupRepository) FindAllByIdLink(ctx context.Context, code string) (info []*entityRam.RamResourceGroupEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id_link like ?", "%|"+code+"|%").Find(&info)
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
