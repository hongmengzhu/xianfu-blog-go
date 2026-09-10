package multiTenantPg

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constContextPg"
	"go-spring.org/log"
	"gorm.io/gorm"
)

// ScopeRulePgWhereOrg 规则
func ScopeRulePgWhereOrg(ctx *gin.Context, tableName string) func(db *gorm.DB) *gorm.DB {
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
		filters, _ := GetMultiOrg(holder)
		//fmt.Printf("typeDomain=%+v\n", holder.GetTypeDomain())
		//fmt.Printf("typeDomain.tableName=%+v\n", tableName)
		//fmt.Printf("typeDomain.multiRule.Tenant=%+\nv", multiRule.Tenant)
		if verifyMultiByTable(&multiTableOrg, tableName) && multiRule.Tenant {
			ruleOrg(db, filters, multiRule)
		}
		return db
	}
}

func ruleOrg(db *gorm.DB, holder MultiTenantPg, rule *MultiRule) *gorm.DB {
	if rule.Org {
		//大于 2个时候
		if nil != holder.OrgNo && len(holder.OrgNo) >= 1 {
			//允许多个 租户列表
			if rule.MultipleOrg && len(holder.OrgNo) > 1 {
				db.Where("org_no IN  ? ", holder.OrgNo)
			} else {
				db.Where("org_no = ?", holder.OrgNo[0])
			}
		} else {
			db.Where("org_no = ?", "000")
		}
	}
	return db
}
