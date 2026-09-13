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
	gs.Provide(new(RamDepartmentRepository))

	gs.Provide(new(support.BaseService[RamDepartmentRepository]))
}

type RamDepartmentRepository struct {
	repositoryPg.BaseRepository[entityRam.RamDepartmentEntity, int64]
}

func (c *RamDepartmentRepository) FindAllByNoLinkArr(ctx context.Context, code []string) (info []*entityRam.RamDepartmentEntity, result bool) {
	db := c.DbModel().WithContext(ctx)
	for index, val := range code {
		if 0 == index {
			db.Where("no_link like ?", "%|"+val+"|%")
		} else {
			db.Or("no_link like ?", "%|"+val+"|%")
		}
	}
	tx := db.Find(&info)
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
