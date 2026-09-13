package repositoryApi

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityApi"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(ApiDiplAccessKeyRepository))

	gs.Provide(new(support.BaseService[ApiDiplAccessKeyRepository]))
}

type ApiDiplAccessKeyRepository struct {
	repositoryPg.BaseRepository[entityApi.ApiDiplAccessKeyEntity, int64]
}

func (c *ApiDiplAccessKeyRepository) FindByTenantNoAndDiplNo(ctx context.Context, no, DiplNo string) (info *entityApi.ApiDiplAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", no).Where("dipl_no=?", DiplNo).First(&info)
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

func (c *ApiDiplAccessKeyRepository) UpdateAllByDiplNoAndNoSetState(ctx context.Context, DiplNo, id string, state int8) (sum int64, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("dipl_no=?", DiplNo).Where("id=?", id).Updates(entityApi.ApiDiplAccessKeyEntity{State: state})
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return 0, false
	}
	if 0 == tx.RowsAffected {
		return 0, false
	}
	return tx.RowsAffected, true
}
