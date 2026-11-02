package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/utilsRam"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcAccount"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsPg/domainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/passwordTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	repositoryRam "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/userPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(NewTcTenantAccountPasswordService)
}

// TcTenantAccountPasswordService 密码修改
// @Description:
type TcTenantAccountPasswordService struct {
	sv    *repositoryRam.RamAccountRepository              `autowire:"?"`
	aAuth *repositoryRam.RamAccountAuthorizationRepository `autowire:"?"`
}

func NewTcTenantAccountPasswordService() *TcTenantAccountPasswordService {
	return new(TcTenantAccountPasswordService)
}

// UpdatePassword 修改密码
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantAccountPasswordService) UpdatePassword(ctx *gin.Context, ct modTcAccount.PasswordCt, option domainPg.Option) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct), log.Any("option", option))
	if ct.ID == 0 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.PasswordNew {
		return rt.ErrorMessage("密码不能为空")
	}
	r := c.sv
	info, b := r.FindByIdAndTypeDomain(ctx, ct.ID.ToInt64(), option.TypeDomain.String())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}

	r2 := c.aAuth
	entity := entityRam.RamAccountAuthorizationEntity{}
	entity.ExtraData = strPg.GetNanoid(8)
	entity.Value = userPg.PasswordSalt(ct.PasswordNew, entity.ExtraData)
	//
	passwd, b := r2.FindByTypePasswordANo(ctx, info.No)
	if !b {
		//密码不存在，则创建
		entity.Ano = info.No
		entity.TenantNo = info.TenantNo
		entity.Type = passwordTypePg.Password.String()
		//设置唯一值
		entity.KindUnique = utilsRam.AuthorizationKindUniquePasswordByEntity(entity)
		r2.Create(ctx, &entity)
	} else {
		r2.Update(ctx, entity, passwd.ID)
	}
	return rt.Ok()
}
