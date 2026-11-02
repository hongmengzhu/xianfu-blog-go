package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamIdpCredentialRepository))

	gs.Provide(new(support.BaseService[RamIdpCredentialRepository]))
}

type RamIdpCredentialRepository struct {
	repositoryPg.BaseRepository[entityRam.RamIdpCredentialEntity, int64]
	//
}

// FindBySourceNoAndCredType 按认证源编号+凭证类型查找
func (c *RamIdpCredentialRepository) FindBySourceNoAndCredType(ctx context.Context, sourceNo, credType string) (info *entityRam.RamIdpCredentialEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("source_no=? AND cred_type=?", sourceNo, credType).First(&info)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}
