package routerPg

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/validatorPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func BindJson(ctx *gin.Context, obj any) bool {
	if err := ctx.ShouldBind(obj); err != nil {
		translate := validatorPg.Translate(err, obj)
		if len(translate) > 0 {
			ctx.JSON(200, rg.ErrorMessageData[string](translate))
			return false
		}
		ctx.JSON(200, rg.ErrorDefault[string]())
		return false
	}
	return true
}
