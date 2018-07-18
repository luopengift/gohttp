package gohttp

import (
	"context"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/luopengift/log"
	"golang.org/x/net/http2"
)

// Application is a httpserver instance.
type Application struct {
	*Config
	*log.Log
	*Template
	*RouterList
	*http.Server
	sync.Pool
}

// Init creates a default httpserver instance by default config.
func Init() *Application {
	app := new(Application)
	app.Config = InitConfig()
	app.Log = InitLog()
	app.Template = InitTemplate()
	app.RouterList = InitRouterList()

	app.Route("^/_routeList$", &RouteHandler{})
	app.Route("^/_info$", &InfoHandler{})

	app.RouteFunc("/debug/pprof/", Index)
	app.RouteFunc("/debug/pprof/cmdline", Cmdline)
	app.RouteFunc("/debug/pprof/profile", Profile)
	app.RouteFunc("/debug/pprof/symbol", Symbol)
	app.RouteFunc("/debug/pprof/trace", Trace)

	app.RouteFunCtx("/debug/gc/start", StartGC)
	app.RouteFunCtx("/debug/gc/stop", StopGC)
	app.RouteFunCtx("/debug/trace/start", StartTrace)
	app.RouteFunCtx("/debug/trace/stop", StopTrace)

	app.Server = &http.Server{
		Addr: app.Config.Addr,
		/** control how to handler ServeHTTP*/
		// Handler:           NewRequestHandler(app),
		Handler:           app,
		ReadTimeout:       time.Duration(app.Config.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(app.Config.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(app.Config.WriteTimeout) * time.Second,
		MaxHeaderBytes:    app.Config.MaxHeaderBytes,
	}
	serverhttp2 := &http2.Server{
		IdleTimeout: 1 * time.Minute,
	}
	if err := http2.ConfigureServer(app.Server, serverhttp2); err != nil {
		app.Error("%v", err)
	}
	app.Pool.New = func() interface{} {
		return &Context{Application: app}
	}
	return app
}

// Run starts the server by listen address.
// HTTP/2.0 is only supported in https,
// If server is http mode, then HTTP/1.x will be used.
func (app *Application) Run(addr ...string) {
	if len(addr) == 1 {
		app.Server.Addr = addr[0]
	} else {
		app.Server.Addr = app.Config.Addr
	}
	app.Info("Http start %s", app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// RunTLS xxx
func (app *Application) RunTLS(addr ...string) {
	if len(addr) == 1 {
		app.Server.Addr = addr[0]
	} else {
		app.Server.Addr = app.Config.Addr
	}
	app.Info("Https start %s", app.Server.Addr)
	if err := app.Server.ListenAndServeTLS(app.Config.CertFile, app.Config.KeyFile); err != nil {
		panic(err)
	}
}

// Stop gracefully shuts down the server without interrupting any active connections.
func (app *Application) Stop() error {
	return app.Server.Shutdown(context.Background())
}

// ServeHTTP is HTTP server implement method. It makes App compatible to native http handler.
func (app *Application) ServeHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	stime := time.Now()
	ctx := app.Pool.Get().(*Context)
	ctx.init(responsewriter, request)
	defer func(ctx *Context) {
		if err := recover(); err != nil {
			debug.PrintStack()
			ctx.HTTPError(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) //500
			ctx.Error(app.LogFormat+" | %v", ctx.Status(), ctx.Method, ctx.URL, ctx.RemoteAddr(), time.Since(stime), err)
		} else {
			switch ctx.Status() / 100 {
			case 2, 3:
				app.Info(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.RemoteAddr(), time.Since(stime))
			case 4:
				app.Warn(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.RemoteAddr(), time.Since(stime))
			case 5:
				app.Error(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.RemoteAddr(), time.Since(stime))
			default:
				app.Error(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.RemoteAddr(), time.Since(stime))
			}
		}
		app.Pool.Put(ctx)

	}(ctx)
	// handler static file
	if strings.HasPrefix(ctx.URL.Path, ctx.Config.StaticPath) || hasSuffixs(ctx.URL.Path, ".ico") {
		file := filepath.Join(ctx.Config.WebPath, ctx.URL.Path)
		http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
		return
	}

	route, match := app.find(ctx.URL.Path)
	if route == nil {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	ctx.match = match
	route.entry.Exec(ctx)
	return
}
