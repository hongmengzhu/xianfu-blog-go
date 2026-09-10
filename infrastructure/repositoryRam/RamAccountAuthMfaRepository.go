package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamAccountAuthMfaRepository))

	gs.Provide(new(support.BaseService[RamAccountAuthMfaRepository]))
}

type RamAccountAuthMfaRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountAuthMfaEntity, int64]
}

// FindByAno 按账号编号查找所有 MFA 凭证
func (c *RamAccountAuthMfaRepository) FindByAno(ctx context.Context, ano string, opts ...optionsPg.Option) (infos []*entityRam.RamAccountAuthMfaEntity, found bool) {
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("ano = ? AND state = 1", ano).Find(&infos)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return infos, true
}

// FindByAnoAndMfaType 按账号编号和MFA类型查找
func (c *RamAccountAuthMfaRepository) FindByAnoAndMfaType(ctx context.Context, ano, mfaType string, opts ...optionsPg.Option) (info *entityRam.RamAccountAuthMfaEntity, found bool) {
	var result entityRam.RamAccountAuthMfaEntity
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("ano = ? AND mfa_type = ? AND state = 1", ano, mfaType).First(&result)
	if tx.Error != nil {
		return nil, false
	}
	return &result, true
}

// FindByAnoAndCredentialID 按账号编号和凭证ID查找
func (c *RamAccountAuthMfaRepository) FindByAnoAndCredentialID(ctx context.Context, ano, credentialID string, opts ...optionsPg.Option) (info *entityRam.RamAccountAuthMfaEntity, found bool) {
	var result entityRam.RamAccountAuthMfaEntity
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("ano = ? AND cred_id = ?", ano, credentialID).First(&result)
	if tx.Error != nil {
		return nil, false
	}
	return &result, true
}

// DeleteByAno 删除账号的所有MFA记录
func (c *RamAccountAuthMfaRepository) DeleteByAno(ctx context.Context, ano string) error {
	return c.DbModel().WithContext(ctx).Where("ano = ?", ano).Delete(&entityRam.RamAccountAuthMfaEntity{}).Error
}

// FindByMfaToken 按MFA临时令牌查找
func (c *RamAccountAuthMfaRepository) FindByMfaToken(ctx context.Context, token string, opts ...optionsPg.Option) (info *entityRam.RamAccountAuthMfaEntity, found bool) {
	var result entityRam.RamAccountAuthMfaEntity
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("mfa_token = ? AND state = 1", token).First(&result)
	if tx.Error != nil {
		return nil, false
	}
	return &result, true
}
