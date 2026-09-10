package repositoryBlog

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogCategoryRepository))

	gs.Provide(new(support.BaseService[BlogCategoryRepository]))
}

type BlogCategoryRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogCategoryEntity, int64]
}

func (c *BlogCategoryRepository) FindAllByParentIdLink(ctx context.Context, code string) (info []*entityBlog.BlogCategoryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id_link like ?", "%"+code+"%").Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *BlogCategoryRepository) FindAllByNoLink(ctx context.Context, code string) (infos []*entityBlog.BlogCategoryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("no_link like ?", "%"+code+"%").Find(&infos)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}
func (c *BlogCategoryRepository) FindAllByCodeLinkAndTypeSys(ctx context.Context, code string, tpSys string) (info []*entityBlog.BlogCategoryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_sys = ?", tpSys).Where("no_link like ?", "%"+code+"%").Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
