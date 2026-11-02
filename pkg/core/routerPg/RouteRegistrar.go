package routerPg

import "github.com/gin-gonic/gin"

type RouteRegistrar interface {
	RegisterRoutes(e *gin.Engine)
}
