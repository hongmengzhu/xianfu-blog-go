package staterEchoV5

import (
	"context"
	"net"
	"net/http"
	"slices"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"go-spring.org/stdlib/errutil"
)

// RouteRegistrar mounts one controller's routes on the engine of the subdomain
// it serves. Every echo-side controller bean is exported as this type, so the
// container collects them into the []RouteRegistrar the server is built from,
// and a controller never has to be named in the wiring.
type RouteRegistrar interface {
	// Host reports the subdomain this controller answers on, for example
	// "blog.test" or "bookmark.test".
	Host() string
	// RegisterRoutes mounts the controller's routes on that subdomain's engine.
	RegisterRoutes(e *echo.Echo)
}

// Server adapts the virtual-host engine to gs.Server: the container starts it
// concurrently with the gin server and shuts it down on SIGINT/SIGTERM.
type Server struct {
	svr   *http.Server
	addr  string
	hosts []string
}

// NewServer builds one Echo engine per subdomain, lets every controller
// registrar mount its routes on the engine of the host it serves, and
// dispatches requests by Host through echo's virtual-host handler. A request
// whose Host matches no subdomain (for example a plain
// http://127.0.0.1:2016/) is served by a fallback engine that lists the
// available subdomains.
func NewServer(regs []RouteRegistrar, cfg Config) (*Server, error) {
	if len(regs) == 0 {
		return nil, errutil.Explain(nil, "echo5: no route registrar registered")
	}

	byHost := make(map[string][]RouteRegistrar, len(regs))
	names := make([]string, 0, len(regs))
	for _, reg := range regs {
		host := reg.Host()
		if host == "" {
			return nil, errutil.Explain(nil, "echo5: route registrar %T has no subdomain", reg)
		}
		if _, ok := byHost[host]; !ok {
			names = append(names, host)
		}
		byHost[host] = append(byHost[host], reg)
	}
	slices.Sort(names)

	// Every subdomain owns its own engine, so two controllers from different
	// subdomains can register the same path without colliding.
	vhosts := make(map[string]*echo.Echo, len(byHost))
	for _, host := range names {
		vh := echo.New()
		vh.Use(middleware.Recover(), middleware.RequestID())
		// 错误处理器：保留 echo 默认响应，同时把原始错误写进日志
		vh.HTTPErrorHandler = newHTTPErrorHandler()
		// 模板渲染器，命名规则与 gin LoadHTMLGlob 保持一致
		if renderer := templatePg.NewEchoRenderer("data/templates/**/**/*"); nil != renderer {
			vh.Renderer = renderer
		}
		// 静态文件目录
		vh.Static("/assets", "./data/assets")
		vh.Static("/attachment", "./data/attachment")
		//
		for _, reg := range byHost[host] {
			reg.RegisterRoutes(vh)
		}
		vhosts[host] = vh
	}

	// echo's virtual-host handler builds its own engine to serve a Host that
	// matches no subdomain, so the fallback routes go on the returned engine:
	// a bare http://127.0.0.1:2016/ then lists the subdomains to try.
	handler := echo.NewVirtualHostHandler(vhosts)
	handler.Use(middleware.Recover(), middleware.RequestID())
	handler.HTTPErrorHandler = newHTTPErrorHandler()
	handler.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"service": "echo5",
			"hosts":   names,
			"hint":    "send Host: blog.test or Host: bookmark.test to reach a subdomain",
		})
	})

	log.Debug(context.Background(), log.TagAppDef, "echo5 virtual hosts prepared", func() []log.Field {
		return []log.Field{log.String("addr", cfg.Address), log.Strings("hosts", names)}
	})

	return &Server{
		addr:  cfg.Address,
		hosts: names,
		svr: &http.Server{
			Addr:              cfg.Address,
			Handler:           hostNormalizer{next: handler},
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
	}, nil
}

// Addr reports the configured listen address (for example ":2016").
func (s *Server) Addr() string { return s.addr }

// Hosts reports the registered subdomains in lexical order.
func (s *Server) Hosts() []string { return s.hosts }

// newHTTPErrorHandler 创建带日志的错误处理器。
// echo v5 的 DefaultHTTPErrorHandler 只负责写响应、不写日志，一旦 handler
// 返回普通 error（模板渲染失败、数据库/服务层错误等），客户端只会看到
// {"message":"Internal Server Error"}，日志里却什么都没有。这里在不改变
// 默认响应的前提下补一条日志：5xx 记 error，4xx 记 debug。
func newHTTPErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		logRequestError(c, err)
		// 响应行为与默认处理器保持一致（不向客户端暴露错误详情）
		echo.DefaultHTTPErrorHandler(false)(c, err)
	}
}

// logRequestError 记录一次失败请求的原始错误及请求上下文。
func logRequestError(c *echo.Context, err error) {
	code := echo.StatusCode(err)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	ctx := context.Background()
	fields := make([]log.Field, 0, 6)
	fields = append(fields, log.Any("status", code))
	if req := c.Request(); nil != req {
		// 带上请求上下文，便于链路追踪
		ctx = req.Context()
		fields = append(fields,
			log.String("method", req.Method),
			log.String("path", req.URL.Path),
			log.String("host", req.Host),
		)
	}
	if rid := c.Response().Header().Get(echo.HeaderXRequestID); "" != rid {
		fields = append(fields, log.String("requestId", rid))
	}
	// 4xx 属于预期内的客户端错误，只记 debug，避免刷屏
	if code < http.StatusInternalServerError {
		log.Debug(ctx, log.TagAppDef, "echo5 request error", func() []log.Field { return fields })
		return
	}
	log.Error(ctx, log.TagAppDef, err, "echo5 request failed", fields...)
}

// hostNormalizer strips the port from the request's Host before echo's
// virtual-host dispatch, so both "Host: blog.test" and "Host:
// blog.test:2016" reach the same engine.
type hostNormalizer struct {
	next http.Handler
}

func (h hostNormalizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if host, _, err := net.SplitHostPort(r.Host); err == nil {
		r.Host = host
	}
	h.next.ServeHTTP(w, r)
}

// Run binds the listener immediately and starts serving after Go-Spring signals
// readiness, so the gin server and this server go live together.
func (s *Server) Run(ctx context.Context, sig gs.ReadySignal) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return errutil.Explain(err, "failed to listen on %s", s.addr)
	}
	<-sig.TriggerAndWait()

	ctx = log.WithFields(ctx, log.String("addr", s.addr), log.Strings("hosts", s.hosts))
	log.Info(ctx, log.TagAppDef, "echo5 server starting")
	err = s.svr.Serve(ln)
	if errutil.IsServerClosed(err) {
		log.Debug(ctx, log.TagAppDef, "echo5 server stopped", nil)
		return nil
	}
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "echo5 server failed")
	}
	return errutil.Explain(err, "failed to serve on %s", s.addr)
}

// Stop gracefully shuts the HTTP server down with the given context, letting
// in-flight requests complete.
func (s *Server) Stop(ctx context.Context) error {
	log.Info(ctx, log.TagAppDef, "echo5 server shutting down", log.String("addr", s.addr))
	return s.svr.Shutdown(ctx)
}
