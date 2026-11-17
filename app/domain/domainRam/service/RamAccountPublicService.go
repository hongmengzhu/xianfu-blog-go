package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modPublic"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccount"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/passwordTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/sessionKeyTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	repositoryRam "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
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

func publicAccountRoles(account holderPg.AccountHolder) []string {
	roles := make([]string, 0, len(account.Os.Roles)+1)
	for _, role := range account.Os.Roles {
		if strPg.IsNotBlank(role) {
			roles = append(roles, role)
		}
	}
	if len(roles) == 0 && account.Founder == 1 {
		roles = append(roles, "administrator")
	}
	return roles
}

// Public 登陆用户信息
func (c *RamAccountPublicService) Public(ctx *gin.Context, holder holderPg.HolderPg) (rt rg.Rs[modPublic.InfoPublicVo]) {
	log.Info(ctx, log.TagAppDef, "holder", log.Any("holder", holder))
	//log.Info(ctx, log.TagAppDef,"HolderData=%+v", holder.HolderData)
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
	data.Info.Roles = publicAccountRoles(account)
	rt.Data = data
	return rt.Ok()
}

// InfoPublic 登陆用户信息
func (c *RamAccountPublicService) InfoPublic(ctx *gin.Context, holder holderPg.HolderPg) (rt rg.Rs[modRamAccount.AccountPub]) {
	log.Info(ctx, log.TagAppDef, "holder", log.Any("holder", holder))
	//log.Info(ctx, log.TagAppDef,"HolderData=%+v", holder.HolderData)
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
	data.Roles = publicAccountRoles(account)
	if len(account.Os.Departments) > 0 {
		data.Departments = account.Os.Departments
	}
	rt.Data = data
	return rt.Ok()
}

// UpdatePassword 修改密码
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamAccountPublicService) UpdatePassword(ctx *gin.Context, ct modPublic.PasswordCt, tp typeDomainPg.TypeDomain) (rt rg.Rs[string]) {
	if "" == ct.PasswordNew {
		return rt.ErrorMessage("密码不能为空")
	}
	if "" == ct.PasswordOld {
		return rt.ErrorMessage("原密码不能为空")
	}
	if len(ct.PasswordNew) < 8 {
		return rt.ErrorMessage("新密码长度不能少于 8 位")
	}
	holder := holderPg.GetContextAccount(ctx)
	account := holder.GetAccount()
	// 密钥对按 域 + 租户 + 客户端 缓存：manage 域用账号租户，system 域用 system 编码
	tenantNo := tp.Code()
	if typeDomainPg.Manage.IsEqual(string(tp)) {
		tenantNo = holder.GetTenantNo()
		if "" == tenantNo || "-1" == tenantNo {
			tenantNo = account.TenantNo
		}
	}
	pwd := ct.PasswordNew
	pwdOld := ct.PasswordOld
	//解密密码：按当前路由域对应的密钥对解密，禁止跨域使用密钥
	if logingEn, ok := c.authLogin.LoginEncrypt["default"]; ok && logingEn {
		login := c.cacheSessionPubPrive.DecodeByLogin(ctx, pwd, tp, clientPg.Browser, tenantNo, sessionKeyTypePg.Login)
		if login.ErrorIs() {
			return rt.ErrorMessage(login.Message)
		}
		pwd = login.Data
		old := c.cacheSessionPubPrive.DecodeByLogin(ctx, pwdOld, tp, clientPg.Browser, tenantNo, sessionKeyTypePg.Login)
		if old.ErrorIs() {
			return rt.ErrorMessage("原密码校验失败")
		}
		pwdOld = old.Data
	}
	if pwd == pwdOld {
		return rt.ErrorMessage("新密码不能与原密码相同")
	}

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
	// 二次确认：原密码必须与库中存储的哈希匹配，防止仅凭 access token 永久改密
	if nil != passwd && !userPg.PasswordVerify(passwd.Value, pwdOld, passwd.ExtraData) {
		return rt.ErrorMessage("原密码错误")
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
