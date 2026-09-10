package multiTenantPg

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constContextPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"go-spring.org/log"
	"gorm.io/gorm"
)

// ScopeRulePgWhere 规则
func ScopeRulePgWhere(ctx *gin.Context, tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var multiRule *MultiRule
		var iHolderRule interfaces.IHolderRule
		value, exists := ctx.Get(constContextPg.CTX_RULE)
		if exists && nil != value {
			log.Errorf(ctx, log.TagAppDef, "CTX_RULE=%+v", value)
			iHolderRule = value.(interfaces.IHolderRule)
			multiRule = iHolderRule.(*MultiRule)
		}
		if nil == multiRule {
			multiRule = &MultiRule{
				MultipleTenant: true,
				Tenant:         true,
				MultipleOrg:    true,
				Merchant:       true,
				MultipleStore:  true,
				Store:          true,
				Owner:          true,
				MultiOwner:     true,
			}
		}
		holder := holderPg.GetContextAccount(ctx)
		//表
		tenant, _ := GetMultiTenant(holder)
		//fmt.Printf("typeDomain=%+v\n", holder.GetTypeDomain())
		//fmt.Printf("typeDomain.tableName=%+v\n", tableName)
		//fmt.Printf("typeDomain.multiRule.Tenant=%+\nv", multiRule.Tenant)
		//系统
		if typeDomainPg.System.IsEqual(holder.GetTypeDomain()) {
			//租户
			if verifyMultiByTable(&multiTableSys, tableName) && multiRule.Tenant {
				ruleTenant(db, tenant, multiRule)
			}
		} else if typeDomainPg.Manage.IsEqual(holder.GetTypeDomain()) {
			// 管理后台
			//租户
			if verifyMultiByTable(&multiTableManage, tableName) && multiRule.Tenant {
				ruleTenant(db, tenant, multiRule)
			}
		} else {
			// 客户
			//租户
			if verifyMultiByTable(&multiTableCustomer, tableName) && multiRule.Tenant {
				ruleTenant(db, tenant, multiRule)
			}
		}

		return db
	}
}

func verifyMultiByTable(multi *pg.MultiItem, tableName string) bool {
	if nil != multi {
		//多租户表，在不包含的表中,直接paas
		if nil != multi.Not && len(multi.Not) > 0 && slice.Contain(multi.Not, tableName) {
			return false
		}
		//多租户表，在不包含的表中,直接paas
		if nil != multi.Contain && len(multi.Contain) > 0 && slice.Contain(multi.Contain, tableName) {
			return true
		}
	}
	return false
}

func ruleTenant(db *gorm.DB, holder MultiTenantPg, rule *MultiRule) *gorm.DB {
	if rule.Tenant {
		//大于 2个时候
		if nil != holder.TenantNo && len(holder.TenantNo) >= 1 {
			//允许多个 租户列表
			if rule.MultipleTenant && len(holder.TenantNo) > 1 {
				db.Where("tenant_no IN  ? ", holder.TenantNo)
			} else {
				db.Where("tenant_no = ?", holder.TenantNo[0])
			}
		} else {
			db.Where("tenant_no = ?", "000")
		}
	}
	return db
}
