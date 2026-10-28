package webPg

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/labstack/echo/v4"
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

// GetTenantNoEcho 从 echo 上下文读取租户编号，与 GetTenantNo 语义一致。
//
//	@Description: 读取租户编号
//	@param c
//	@return string
func GetTenantNoEcho(c echo.Context) string {
	if value := c.Get(constHeaderPg.WebTenantNo); nil != value {
		if no, ok := value.(string); ok {
			return no
		}
	}
	return "-1"
}
