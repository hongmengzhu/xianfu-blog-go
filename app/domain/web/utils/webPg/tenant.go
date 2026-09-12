package webPg

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
)

// GetTenantNo
//
//	@Description:
//	@param ctx
//	@return string
func GetTenantNo(ctx *gin.Context) string {
	value, exists := ctx.Get(constHeaderPg.WebTenantNo)
	if exists {
		return value.(string)
	}
	return "-1"
}
