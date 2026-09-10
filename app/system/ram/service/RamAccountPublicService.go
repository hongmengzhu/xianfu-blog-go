package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheRam"
	modPublic2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modPublic"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccount"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	holderPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/passwordTypePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/userPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(NewRamAccountPublicService)
}

// RamAccountPublicService 账户公共动作
// @Description:
type RamAccountPublicService struct {
	sv                   *repositoryRam.RamAccountRepository              `autowire:"?"`
	aAuth                *repositoryRam.RamAccountAuthorizationRepository `autowire:"?"`
	authLogin            pg.Auth                                          `value:"${pg.auth}"`
	cacheSessionPubPrive *cacheRam.CacheSessionPubPrive                   `autowire:"?" `
}

func NewRamAccountPublicService() *RamAccountPublicService {
	return new(RamAccountPublicService)
}

// Public 登陆用户信息
func (c *RamAccountPublicService) Public(ctx *gin.Context, holder holderPg2.HolderPg) (rt rg.Rs[modRamAccount.AccountPub]) {
	log.Infof(ctx, log.TagAppDef, "holder=%+v", holder)
	log.Infof(ctx, log.TagAppDef, "HolderData=%+v", holder.HolderData)
	if nil == holder.HolderData {
		return rt.ErrorMessage("账号登陆失败")
	}
	data := rt.Data
	account := holder.GetAccount()
	copier.Copy(&data, &account)
	data.RealName = account.Name
	data.Avatar = ""
	data.Username = account.Account
	data.UserId = numberPg.Int64ToString(account.ID)
	data.Departments = make([]string, 0)
	if len(account.Os.Departments) > 0 {
		data.Departments = account.Os.Departments
	}
	rt.Data = data
	return rt.Ok()
}

// InfoPublic 登陆用户信息
func (c *RamAccountPublicService) InfoPublic(ctx *gin.Context, holder holderPg2.HolderPg) (rt rg.Rs[modPublic2.InfoPublicVo]) {
	log.Infof(ctx, log.TagAppDef, "holder=%+v", holder)
	log.Infof(ctx, log.TagAppDef, "HolderData=%+v", holder.HolderData)
	if nil == holder.HolderData {
		return rt.ErrorMessage("账号登陆失败")
	}
	data := rt.Data
	account := holder.GetAccount()
	copier.Copy(&data.Info, &account)
	data.Info.RealName = account.Name
	data.Info.Avatar = ""
	data.Info.Username = account.Account
	data.Info.UserId = numberPg.Int64ToString(account.ID)
	data.Info.Departments = make([]string, 0)
	if len(account.Os.Departments) > 0 {
		data.Info.Departments = account.Os.Departments
	}
	data.Info.Roles = make([]string, 0)
	data.Info.Roles = append(data.Info.Roles, "administrator")
	rt.Data = data
	return rt.Ok()
}

// UpdatePassword 修改密码
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamAccountPublicService) UpdatePassword(ctx *gin.Context, ct modPublic2.PasswordCt) (rt rg.Rs[string]) {
	if "" == ct.PasswordNew {
		return rt.ErrorMessage("密码不能为空")
	}
	pwd := ct.PasswordNew
	//解密密码
	if logingEn, ok := c.authLogin.LoginEncrypt["default"]; ok && logingEn {
		login := c.cacheSessionPubPrive.DecodeByLoginSystem(ctx, pwd)
		if login.ErrorIs() {
			return rt.ErrorMessage(login.Message)
		}
		pwd = login.Data
	}

	holder := holderPg2.GetContextAccount(ctx)
	account := holder.GetAccount()
	r := c.sv
	info, b := r.FindById(ctx, account.ID)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	r2 := c.aAuth
	passwd, result := r2.FindByTypePasswordANo(ctx, info.No)
	if !result {
		return rt.ErrorMessage("数据不存在")
	}
	entity := entityRam.RamAccountAuthorizationEntity{}
	entity.ExtraData = strPg.GetNanoid(8)
	entity.Value = userPg.PasswordSalt(pwd, entity.ExtraData)
	if nil == passwd {
		entity.Ano = info.No
		entity.TenantNo = info.TenantNo
		entity.Type = passwordTypePg.Password.String()
		r2.Create(ctx, &entity)
	} else {
		r2.Update(ctx, entity, passwd.ID)
	}

	return rt.Ok()
}
