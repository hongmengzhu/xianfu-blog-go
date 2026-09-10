package repositoryBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicCountryRepository))

	gs.Provide(new(support.BaseService[BasicCountryRepository]))
}

type BasicCountryRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicCountryEntity, int64]
}

func (c *BasicCountryRepository) FindByCountryCode(ctx context.Context, code string) (info *entityBasic.BasicCountryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("state=1").Where("country_code=?", code).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *BasicCountryRepository) FindByCountryCodeAndIdNot(ctx context.Context, code, id string) (info *entityBasic.BasicCountryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("state=1").Where("country_code=?", code).Where("id!=?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
