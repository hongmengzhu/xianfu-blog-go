package multiTenantPg

import (
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"gorm.io/gorm"
)

// ScopeRulePgWhere 规则
func ScopeRulePgWhere(ctx *gin.Context, tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		multiRule := resolveMultiRule(ctx)
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
			// 系统域是平台超管，按设计需要跨租户管理数据，不做租户过滤。
			return db
		} else if typeDomainPg.Manage.IsEqual(holder.GetTypeDomain()) {
			// 管理后台
			//租户
			if verifyMultiByTable(&multiTableManage, tableName) && multiRule.Tenant {
				db = ruleTenant(db, tenant, multiRule)
			}
		} else if typeDomainPg.Customer.IsEqual(holder.GetTypeDomain()) {
			// 客户
			//租户
			if verifyMultiByTable(&multiTableCustomer, tableName) && multiRule.Tenant {
				db = ruleTenant(db, tenant, multiRule)
			}
		} else {
			// 客户
			//租户
			if verifyMultiByTable(&multiTableCustomer, tableName) && multiRule.Tenant {
				db = ruleTenant(db, tenant, multiRule)
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
	if !rule.Tenant {
		return db
	}
	// 大于 1 个时候允许租户列表
	if len(holder.TenantNo) >= 1 {
		if rule.MultipleTenant && len(holder.TenantNo) > 1 {
			return db.Where("tenant_no IN ?", holder.TenantNo)
		}
		return db.Where("tenant_no = ?", holder.TenantNo[0])
	}
	// 没有租户上下文时 fail-closed，避免退化为全表查询
	return db.Where("tenant_no = ?", "000")
}
