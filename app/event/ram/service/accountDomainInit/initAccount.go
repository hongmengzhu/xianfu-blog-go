package accountDomainInit

import (
	"context"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/utilsRam"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/dataFilePg"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/authorizationTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeIdentityPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/userPg"
	"go-spring.org/log"
	"gorm.io/datatypes"
)

var domain = []string{
	typeDomainPg.System.Code(),
	//iamConstant.Merchant.Index(),
	//iamConstant.Tenant.Index(),
	typeDomainPg.Manage.Code(),
}

// InitAccount
// @Description: 初始化账号
type InitAccount struct {
	sp *Sp `autowire:"?"`
}

func NewInitAccount(sp *Sp) *InitAccount {
	return &InitAccount{
		sp: sp,
	}
}

func (t *InitAccount) Processor(ctx context.Context) error {
	mapDomain := make(map[string]*entityRam.RamAccountEntity)
	find := make(map[string]bool)
	// 获取超管账号
	info, result := t.sp.accDb.FindByTypeDomainIn(ctx, domain)
	if result {
		if len(info) > 0 {
			for _, v := range info {
				mapDomain[v.No] = v
				//系统
				if typeDomainPg.System.IsEqual(v.No) {
					find[v.No] = true
				}
				//租户
				if constsPg.ACCOUNT_MANAGE_No == v.No {
					find[typeDomainPg.Manage.Code()] = true
				}
			}
		}
	}
	for _, item := range domain {
		//不存在时创建
		if _, ok := find[item]; !ok {
			//系统账号
			if typeDomainPg.System.IsEqual(item) {
				t.systemAccount(ctx, item)
			} else if typeDomainPg.Manage.IsEqual(item) {
				//管理后台账号
				t.manageAccount(ctx, item)
			}
		} else {
			log.Debugf(ctx, log.TagAppDef, "[init].[账号初始化] %+v 已存在,无需初始化", item)
		}
	}
	return nil
}

func (t *InitAccount) systemAccount(ctx context.Context, domain string) {
	no := noPg.No()
	now := time.Now()
	save := entityRam.RamAccountEntity{
		ID:            1,
		No:            domain,
		Code:          domain,
		TenantNo:      domain,
		TypeDomain:    domain,
		TypeIdentity:  typeIdentityPg.General.Index(),
		Account:       constsPg.ACCOUNT_SYSTEM,
		Cc:            constsPg.ACCOUNT_CC,
		Mail:          no + "@xxxx.com",
		AccountVerify: enumStatePg.ENABLE.Index(),
		State:         enumStatePg.ENABLE.Index(),
		Founder:       enumStatePg.ENABLE.Index(),
		RegisterTime:  &now,
	}
	save.AccountMd5 = cryptPg.Md5(save.Account)
	save.PhoneMd5 = cryptPg.Md5(save.Account)
	save.MailMd5 = cryptPg.Md5(save.Account)
	save.Name = constsPg.ACCOUNT_NAME
	save.RealName = constsPg.ACCOUNT_NAME
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
	os.Tenants = append(os.Tenants, save.TenantNo)
	save.Os = datatypes.NewJSONType(os)
	t.sp.accDb.Create(ctx, &save)
	//
	t.sp.authDb.DeleteByAno(ctx, save.No)
	//
	salt := strPg.GetNanoid(8)
	authorizationEntity := entityRam.RamAccountAuthorizationEntity{
		Ano:       save.No,
		Type:      authorizationTypePg.PASSWORD.String(),
		TenantNo:  save.TenantNo,
		ExtraData: salt,
	}
	pwd := constsPg.ACCOUNT_PASSWORD + "." + strPg.GetNanoid(8)
	authorizationEntity.Value = userPg.PasswordSalt(pwd, salt)
	//设置唯一值
	authorizationEntity.KindUnique = utilsRam.AuthorizationKindUniquePasswordByEntity(authorizationEntity)
	t.sp.authDb.Create(ctx, &authorizationEntity)
	//
	//记录文件日志
	dataFilePg.NewAccountFileRecord(t.sp.pg, dataFilePg.MakeContent(save.Account, pwd, typeDomainPg.System.Code())).Write()
}

func (t *InitAccount) manageAccount(ctx context.Context, domain string) {
	no := noPg.No()
	now := time.Now()
	tenant, result := t.sp.tenantDb.FindByNo(ctx, constsPg.ACCOUNT_MANAGE_No)
	if !result {
		tenant = &entityTc.TcTenantEntity{
			ID:       1000,
			No:       constsPg.ACCOUNT_MANAGE_No,
			Name:     "默认",
			NameFl:   "默认",
			NameFull: "默认",
			Code:     constsPg.ACCOUNT_MANAGE_No,
			State:    enumStatePg.ENABLE.Index(),
			Founder:  constsPg.ACCOUNT_MANAGE_No,
		}
		t.sp.tenantDb.Create(ctx, tenant)
		//
	}
	save := entityRam.RamAccountEntity{
		ID:            1000,
		No:            constsPg.ACCOUNT_MANAGE_No,
		Code:          constsPg.ACCOUNT_MANAGE_No,
		TenantNo:      tenant.No,
		TypeDomain:    domain,
		TypeIdentity:  typeIdentityPg.General.Index(),
		Account:       constsPg.ACCOUNT_MANAGE,
		Cc:            constsPg.ACCOUNT_CC,
		Mail:          no + "@xxxx.com",
		AccountVerify: enumStatePg.ENABLE.Index(),
		State:         enumStatePg.ENABLE.Index(),
		Founder:       enumStatePg.ENABLE.Index(),
		RegisterTime:  &now,
	}
	save.AccountMd5 = cryptPg.Md5(save.Account)
	save.PhoneMd5 = cryptPg.Md5(save.Account)
	save.MailMd5 = cryptPg.Md5(save.Account)
	save.Name = constsPg.ACCOUNT_NAME
	save.RealName = constsPg.ACCOUNT_NAME
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
	os.Tenants = append(os.Tenants, save.TenantNo)
	save.Os = datatypes.NewJSONType(os)
	t.sp.accDb.Create(ctx, &save)
	//
	t.sp.authDb.DeleteByAno(ctx, save.No)
	//
	salt := strPg.GetNanoid(8)
	authorizationEntity := entityRam.RamAccountAuthorizationEntity{
		Ano:       save.No,
		Type:      authorizationTypePg.PASSWORD.String(),
		TenantNo:  save.TenantNo,
		ExtraData: salt,
	}
	authorizationEntity.Value = userPg.PasswordSalt(constsPg.ACCOUNT_PASSWORD, salt)
	//设置唯一值
	authorizationEntity.KindUnique = utilsRam.AuthorizationKindUniquePasswordByEntity(authorizationEntity)
	t.sp.authDb.Create(ctx, &authorizationEntity)
}
