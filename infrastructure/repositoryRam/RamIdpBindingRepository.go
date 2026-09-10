package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamIdpBindingRepository))

	gs.Provide(new(support.BaseService[RamIdpBindingRepository]))
}

type RamIdpBindingRepository struct {
	repositoryPg.BaseRepository[entityRam.RamIdpBindingEntity, int64]
	//
}

// FindByIdpAndExternalSub 按提供商+外部身份标识查找绑定
func (c *RamIdpBindingRepository) FindByIdpAndExternalSub(ctx context.Context, idp, externalSub string) (info *entityRam.RamIdpBindingEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("idp=? AND external_sub=?", idp, externalSub).First(&info)
	if tx.Error != nil {
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindByIdpAndOpenId 按提供商+OpenId查找绑定
func (c *RamIdpBindingRepository) FindByIdpAndOpenId(ctx context.Context, idp, openId string) (info *entityRam.RamIdpBindingEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("idp=? AND open_id=?", idp, openId).First(&info)
	if tx.Error != nil {
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindByIdpAndUnionId 按提供商+UnionId查找绑定
func (c *RamIdpBindingRepository) FindByIdpAndUnionId(ctx context.Context, idp, unionId string) (info *entityRam.RamIdpBindingEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("idp=? AND union_id=?", idp, unionId).First(&info)
	if tx.Error != nil {
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}
