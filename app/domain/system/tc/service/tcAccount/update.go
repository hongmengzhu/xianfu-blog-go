package tcAccount

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcAccount"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/identityPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/sexPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/appModulePg"
	"go-spring.org/log"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"gorm.io/datatypes"
)

type Update struct {
	sp  *Sp `autowire:"?"`
	ctx *gin.Context
	//
	entity *entityRam.RamAccountEntity
}

// NewUpdate
//
//	@Description: 更新
//	@param log
//	@param accDb
//	@param authDb
//	@param depDb
//	@param roleDb
//	@param teamDb
//	@param ctx
//	@return *Update
func NewUpdate(
	sp *Sp,
	ctx *gin.Context) *Update {
	return &Update{
		sp:     sp,
		ctx:    ctx,
		entity: &entityRam.RamAccountEntity{},
	}
}

// accountUpdate 更新
//
//	@Description:
//	@receiver c
//	@param ct
func (c *Update) accountUpdate(ctx *gin.Context, ct modTcAccount.UpdateAccountCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if ct.ID <= 0 {
		return rt.ErrorMessage("id不能为空")
	}
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
	query := false
	c.entity, query = r.FindByIdAndTypeDomain(ctx, ct.ID.ToInt64(), tp.ToTypeDomain().String())
	if !query {
		return rt.ErrorMessage("账号信息不存在")
	}
	find, ok := r.FindByAccountAndTypeDomain(ctx, ct.Account, tp.ToTypeDomain().String())
	if ok && ct.ID.ToInt64() != find.ID {
		return rt.ErrorMessage("账户已存在")
	}
	_, ok = r.FindByPhoneAndTypeDomainAndIdNot(ctx, ct.Phone, tp.ToTypeDomain().String(), ct.ID.ToString())
	if ok {
		return rt.ErrorMessage("手机号已存在")
	}
	_, ok = r.FindByMailAndTypeDomainAndIdNot(ctx, ct.Mail, tp.ToTypeDomain().String(), ct.ID.ToString())
	if ok {
		return rt.ErrorMessage("邮箱已存在")
	}
	_, ok = r.FindByNoAndTypeDomainAndIdNot(ctx, ct.Code, tp.ToTypeDomain().String(), ct.ID.ToString())
	if ok {
		return rt.ErrorMessage("编号已存在")
	}
	//
	var entity entityRam.RamAccountEntity
	copier.Copy(&entity, &ct)
	//
	entity.TenantNo = ""
	//
	entity.MailMd5 = cryptPg.Md5(entity.Mail)
	entity.PhoneMd5 = cryptPg.Md5(entity.Phone)
	entity.AccountMd5 = cryptPg.Md5(entity.Account)
	//
	entity.No = ""
	log.Infof(ctx, log.TagAppDef, "update=%+v", entity)
	err := r.Update(c.ctx, entity, entity.ID)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "save.error=%#v", err)
		return rt.ErrorMessage("保存失败")
	}
	return rt.Ok()
}

// UpdateAccount 更新
//
//	@Description:
//	@receiver c
//	@param ct
func (c *Update) UpdateAccount(ctx *gin.Context, ct modTcAccount.UpdateAccountCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	account := c.accountUpdate(ctx, ct, tp)
	if account.ErrorIs() {
		return rt.ErrorMessage(account.Message)
	}
	return rt.Ok()
}

// updateAll
//
//	@Description: 更新所有
//	@receiver c
//	@param ctx
//	@param ct
//	@param tp
//	@return rt
func (c *Update) updateAll(ctx *gin.Context, ct modTcAccount.UpdateCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var ctAccount modTcAccount.UpdateAccountCt
	copier.Copy(&ctAccount, &ct)
	account := c.accountUpdate(ctx, ctAccount, tp)
	if account.ErrorIs() {
		return rt.ErrorMessage(account.Message)
	}
	//转换
	var dataCt modTcAccount.UpdateDataCt
	copier.Copy(&dataCt, &ct)
	//赋值
	var entity entityRam.RamAccountEntity
	copier.Copy(&entity, &dataCt)
	//
	r := c.sp.accDb
	info, query := r.FindByIdAndTypeDomain(ctx, ct.ID.ToInt64(), tp.ToTypeDomain().String())
	if !query {
		return rt.ErrorMessage("账号信息不存在")
	}
	os := info.Os.Data()
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
	entity.ID = 0
	entity.No = ""
	log.Infof(ctx, log.TagAppDef, "update=%+v", entity)
	err := r.Update(c.ctx, entity, info.ID)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "save.error=%#v", err)
		return rt.ErrorMessage("保存失败")
	}
	return rt.OkData(numberPg.Int64ToString(info.ID))
}

// Process
//
//	@Description: 处理
//	@receiver c
//	@param ctx
//	@param ct
//	@param tp
//	@return rt
func (c *Update) Process(ctx *gin.Context, ct modTcAccount.UpdateCt, tp appModulePg.AppModule) (rt rg.Rs[string]) {
	return c.updateAll(ctx, ct, tp)
}
