package templatePg

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type HttpPg struct {
	Path string   `json:"path"`
	URL  *url.URL `json:"URL"`
}

func NewHttpPg(ctx *gin.Context) *HttpPg {
	return NewHttpPgRequest(ctx.Request)
}

// NewHttpPgRequest 由 *http.Request 构造，供 echo 等其它框架复用。
func NewHttpPgRequest(req *http.Request) *HttpPg {
	hp := &HttpPg{}
	if nil != req && nil != req.URL {
		hp.Path = req.URL.Path
		hp.URL = req.URL
	}
	return hp
}
