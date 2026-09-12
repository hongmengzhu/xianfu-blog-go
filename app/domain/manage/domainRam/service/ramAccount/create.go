package ramAccount

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/utilsRam"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccount"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/authorizationTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/identityPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/sexPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/appModulePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"go-spring.org/log"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"gorm.io/datatypes"
)

type Create struct {
	sp  *Sp `autowire:"?"`
	ctx *gin.Context
	//
	entity *entityRam.RamAccountEntity
}

// NewCreate
//
//	@Description: 创建
//	@param log
//	@param accDb
//	@param authDb
//	@param depDb
//	@param roleDb
//	@param teamDb
//	@param ctx
//	@return *Create
func NewCreate(
	sp *Sp,
	ctx *gin.Context) *Create {
	return &Create{
		sp:     sp,
		ctx:    ctx,
		entity: &entityRam.RamAccountEntity{},
	}
}

// accountCreateSimple 创建
//
//	@Description:
//	@receiver c
//	@param ct
func (c *Create) accountCreateSimple(ctx *gin.Context, ct modRamAccount.CreateUpdateAccountCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if len(ct.Account) <= 0 {
		return rt.ErrorMessage("账户不能为空")
	}
	if len(ct.Phone) <= 0 {
		return rt.ErrorMessage("手机号不能为空")
	}
	if len(ct.Mail) <= 0 {
		return rt.ErrorMessage("邮箱不能为空")
	}
	if len(ct.Code) <= 0 {
		ct.Code = strPg.GetNanoid(10)
	}
	r := c.sp.accDb
	_, ok := r.FindByAccountAndTypeDomainAndIdNot(ctx, ct.Account, tp.ToTypeDomain().String(), "0")
	if ok {
		return rt.ErrorMessage("账户已存在")
	}
	_, ok = r.FindByPhoneAndTypeDomainAndIdNot(ctx, ct.Phone, tp.ToTypeDomain().String(), "0")
	if ok {
		return rt.ErrorMessage("手机号已存在")
	}
	_, ok = r.FindByMailAndTypeDomainAndIdNot(ctx, ct.Mail, tp.ToTypeDomain().String(), "0")
	if ok {
		return rt.ErrorMessage("邮箱已存在")
	}
	_, ok = r.FindByNoAndTypeDomainAndIdNot(ctx, ct.Code, tp.ToTypeDomain().String(), "0")
	if ok {
		return rt.ErrorMessage("编号已存在")
	}
	copier.Copy(&c.entity, &ct)
	c.entity.ID = 0
	c.entity.ExtraData = datatypes.JSON([]byte(`{}`))
	c.entity.TypeDomain = tp.ToTypeDomain().String() //域类型
	//
	c.entity.Name = c.entity.Account
	c.entity.RealName = c.entity.Account
	if strPg.IsBlank(c.entity.Sex) {
		c.entity.Sex = sexPg.MALE.String()
	}
	//
	os := entityRam.RamAccountJsonOs{
		Departments: make([]string, 0),
		Groups:      make([]string, 0),
		Levels:      make([]string, 0),
		Merchants:   make([]string, 0),
		Orgs:        make([]string, 0),
		Projects:    make([]string, 0),
		Roles:       make([]string, 0),
		Shops:       make([]string, 0),
		Teams:       make([]string, 0),
		Tenants:     make([]string, 0),
	}
	//
	c.entity.Os = datatypes.NewJSONType(os)
	// 国际
	if strPg.IsBlank(ct.Cc) {
		ct.Cc = "86"
	}
	now := time.Now()
	holder := holderPg.GetContextAccount(c.ctx)
	//
	c.entity.MailMd5 = cryptPg.Md5(c.entity.Mail)
	c.entity.PhoneMd5 = cryptPg.Md5(c.entity.Phone)
	c.entity.AccountMd5 = cryptPg.Md5(c.entity.Account)
	if nil == c.entity.RegisterTime {
		c.entity.RegisterTime = &now
	}
	c.entity.No = noPg.No()
	c.entity.TenantNo = holder.GetTenantNo()
	log.Infof(ctx, log.TagAppDef, "save=%#v", c.entity)
	err, _ := r.Create(c.ctx, c.entity)
	if err != nil {
		return rt.ErrorMessage("创建用户失败")
	}
	authorizationEntity := entityRam.RamAccountAuthorizationEntity{
		Ano:      c.entity.No,
		Type:     authorizationTypePg.PASSWORD.String(),
		TenantNo: c.entity.TenantNo,
	}
	//设置唯一值
	authorizationEntity.KindUnique = utilsRam.AuthorizationKindUniquePasswordByEntity(authorizationEntity)
	c.sp.authDb.Create(c.ctx, &authorizationEntity)
	return rt.Ok()
}

// CreateAccountSimple 创建
//
//	@Description:
//	@receiver c
//	@param ct
func (c *Create) CreateAccountSimple(ctx *gin.Context, ct modRamAccount.CreateUpdateAccountCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	account := c.accountCreateSimple(ctx, ct, tp)
	if account.ErrorIs() {
		return rt.ErrorMessage(account.Message)
	}
	return rt.OkData(numberPg.Int64ToString(c.entity.ID))
}

// createAll
//
//	@Description: 创建所有
//	@receiver c
//	@param ctx
//	@param ct
//	@param tp
//	@return rt
func (c *Create) createAll(ctx *gin.Context, ct modRamAccount.CreateUpdateCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var ctAccount modRamAccount.CreateUpdateAccountCt
	copier.Copy(&ctAccount, &ct)
	ctAccount.ID = 0
	account := c.accountCreateSimple(ctx, ctAccount, tp)
	if account.ErrorIs() {
		return rt.ErrorMessage(account.Message)
	}
	var dataCt modRamAccount.CreateUpdateCt
	copier.Copy(&dataCt, &ct)
	var entity entityRam.RamAccountEntity
	//赋值 数据
	copier.Copy(&entity, &dataCt)
	//
	//身份
	{
		if strPg.IsBlank(entity.TypeIdentity) {
			entity.TypeIdentity = identityPg.GENERAL.String()
		}
		if !identityPg.IsExist(entity.TypeIdentity) {
			return rt.ErrorMessage("身份错误")
		}
	}
	//性别
	{
		if strPg.IsBlank(entity.Sex) {
			entity.Sex = sexPg.MALE.String()
		}
		if !sexPg.IsExist(entity.Sex) {
			return rt.ErrorMessage("性别错误")
		}
	}
	//
	os := entity.Os.Data()
	os.Departments = make([]string, 0)
	os.Roles = make([]string, 0)
	os.Levels = make([]string, 0)
	os.Groups = make([]string, 0)
	os.Teams = make([]string, 0)
	//多部门
	if nil != ct.Departments && len(ct.Departments) > 0 {
		ids := make([]string, 0)
		for _, id := range ct.Departments {
			if "" != id {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			infos, b := c.sp.depDb.FindAllByNoIn(c.ctx, ids)
			if b {
				for _, item := range infos {
					os.Departments = append(os.Departments, item.No)
				}
				os.Departments = ids
			} else {
				return rt.ErrorMessage("部分部门匹配失败")
			}
		} else {
			return rt.ErrorMessage("部分部门匹配失败")
		}
	}
	//多角色
	if nil != ct.Roles && len(ct.Roles) > 0 {
		ids := make([]string, 0)
		for _, id := range ct.Roles {
			if strPg.IsNotBlank(id) {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			infos, b := c.sp.roleDb.FindAllByNoIn(c.ctx, ids)
			if b {
				for _, item := range infos {
					os.Roles = append(os.Roles, item.No)
				}
			} else {
				return rt.ErrorMessage("部分角色匹配失败")
			}
		} else {
			return rt.ErrorMessage("部分角色匹配失败")
		}
	}
	//多级别
	if nil != ct.Levels && len(ct.Levels) > 0 {
		ids := make([]string, 0)
		for _, id := range ct.Levels {
			if "" != id {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			infos, b := c.sp.levelDb.FindAllByNoIn(c.ctx, ids)
			if b {
				for _, item := range infos {
					os.Levels = append(os.Levels, item.No)
				}
			} else {
				return rt.ErrorMessage("部分级别匹配失败")
			}
		} else {
			return rt.ErrorMessage("部分级别匹配失败")
		}
	}
	//多组
	if nil != ct.Groups && len(ct.Groups) > 0 {
		ids := make([]string, 0)
		for _, id := range ct.Groups {
			if "" != id {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			infos, b := c.sp.groupDb.FindAllByNoIn(c.ctx, ids)
			if b {
				for _, item := range infos {
					os.Groups = append(os.Groups, item.No)
				}
			} else {
				return rt.ErrorMessage("部分分组匹配失败")
			}
		} else {
			return rt.ErrorMessage("部分分组匹配失败")
		}
	}
	//多团队
	if nil != ct.Teams && len(ct.Teams) > 0 {
		ids := make([]string, 0)
		for _, id := range ct.Teams {
			if "" != id {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			infos, b := c.sp.teamDb.FindAllByNoIn(c.ctx, ids)
			if b {
				for _, item := range infos {
					os.Teams = append(os.Teams, item.No)
				}
			} else {
				return rt.ErrorMessage("部分团队匹配失败")
			}
		} else {
			return rt.ErrorMessage("部分团队匹配失败")
		}
	}
	//
	entity.Os = datatypes.NewJSONType(os)
	//
	log.Infof(ctx, log.TagAppDef, "save=%#v", entity)
	err := c.sp.accDb.Update(c.ctx, entity, c.entity.ID)
	if err != nil {
		return rt.ErrorMessage("创建用户失败")
	}
	return rt.OkData(numberPg.Int64ToString(entity.ID))
}

// Process
//
//	@Description: 处理
//	@receiver c
//	@param ctx
//	@param ct
//	@param tp
//	@return rt
func (c *Create) Process(ctx *gin.Context, ct modRamAccount.CreateUpdateCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	return c.createAll(ctx, ct, tp)
}
