package templatePg

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

type HttpPg struct {
	Path string   `json:"path"`
	URL  *url.URL `json:"URL"`
}

func NewHttpPg(ctx *gin.Context) *HttpPg {
	return &HttpPg{Path: ctx.Request.URL.Path, URL: ctx.Request.URL}
}
