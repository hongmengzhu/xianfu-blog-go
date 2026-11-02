package holderApiPg

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constContextPg"
)

func GetContextAccount(ctx *gin.Context) (holder HolderPg) {
	get, is := ctx.Get(constContextPg.AUTH_LOGIN_API)
	if !is || nil == get {
		return
	}
	holder = get.(HolderPg)
	return
}

func SetContextValue(ctx *gin.Context) context.Context {
	return context.WithValue(ctx.Request.Context(), constContextPg.CTX, ctx)
}

func SetContextValueGs(ctx *gin.Context) context.Context {
	return context.WithValue(ctx.Request.Context(), constContextPg.CTX, ctx)
}
