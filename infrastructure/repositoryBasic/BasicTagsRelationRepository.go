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
	gs.Provide(new(BasicTagsRelationRepository))

	gs.Provide(new(support.BaseService[BasicTagsRelationRepository]))
}

type BasicTagsRelationRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicTagsRelationEntity, int64]
}

func (c *BasicTagsRelationRepository) FindByNameAndIdNotAndCategoryNot(ctx context.Context, name string, id int64, category string) (info *entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("name=?", name).Where("category=?", category).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicTagsRelationRepository) FindByCodeAndIdNotAndCategoryNot(ctx context.Context, name string, id int64, category string) (info *entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("code=?", name).Where("category=?", category).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (b *BasicTagsRelationRepository) DeleteByIdsStringAndTypeSysNot(ctx context.Context, id []string, tp string) error {
	tx := b.DbModel().WithContext(ctx).Where("type_sys!=?", tp).Where("id in ?").Delete(&b.Entity, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// FindAllByCategoryNoIn
//
//	@Description:
//	@receiver c
//	@param t
//	@param category
//	@return infos
//	@return result
func (c *BasicTagsRelationRepository) FindAllByCategoryNoIn(ctx context.Context, t entityBasic.BasicTagsRelationEntity, category []string) (infos []*entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where(t).Where("category_no in ?", category).Find(&infos)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

// FindAllByCategoryRootIn
//
//	@Description:
//	@receiver c
//	@param t
//	@param category
//	@return infos
//	@return result
func (c *BasicTagsRelationRepository) FindAllByCategoryRootIn(ctx context.Context, t entityBasic.BasicTagsRelationEntity, category []string) (infos []*entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where(t).Where("category_root in ?", category).Find(&infos)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

// FindAllByCodeInAndCategoryRoot
//
//	@Description:
//	@receiver c
//	@param t
//	@param category
//	@return infos
//	@return result
func (c *BasicTagsRelationRepository) FindAllByCodeInAndCategoryRoot(ctx context.Context, code []string, categoryRoot string) (infos []*entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("code in ?", code).Where("category_root=?", categoryRoot).Find(&infos)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

// FindAllByTagNoInAndCategoryRoot
//
//	@Description:
//	@receiver c
//	@param t
//	@param category
//	@return infos
//	@return result
func (c *BasicTagsRelationRepository) FindAllByTagNoInAndCategoryRoot(ctx context.Context, code []string, categoryRoot string) (infos []*entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tag_no in ?", code).Where("category_root=?", categoryRoot).Find(&infos)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}
