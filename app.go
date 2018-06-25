/*
   gohttp sample http server framework
*/
package gohttp

import (
	"context"
	"github.com/luopengift/log"
	"golang.org/x/net/http2"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

// Aplllication is a httpserver instance.
type Application struct {
	*Config
	*log.Log
	*Template
	*RouterList
	*http.Server
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
	serverHttp2 := &http2.Server{
		IdleTimeout: 1 * time.Minute,
	}
	if err := http2.ConfigureServer(app.Server, serverHttp2); err != nil {
		app.Error("%v", err)
	}
	return app
}

// Run starts the server by listen address.
// HTTP/2.0 is only supported in https,
// If server is http mode, then HTTP/1.x will be used.
func (app *Application) RunHttp(addr ...string) {
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

func (app *Application) RunHttps(addr ...string) {
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
	// init a new http handler
	ctx := NewHttpHandler(app, responsewriter, request)

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			ctx.HTTPError(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) //500
			ctx.Error(app.LogFormat+" | %v", ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime), err)
		}
	}()
	app.handler(ctx)
	switch ctx.Status() / 100 {
	case 2, 3:
		ctx.Info(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	case 4:
		ctx.Warn(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	case 5:
		ctx.Error(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	default:
		ctx.Error(app.LogFormat, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	}
}

//func (app *Application) handler(responsewriter http.ResponseWriter, request *http.Request) {
func (app *Application) handler(ctx *HttpHandler) {
	if strings.HasPrefix(ctx.Path, ctx.Config.StaticPath) || hasSuffixs(ctx.Path, ".ico") {
		file := filepath.Join(ctx.Config.WebPath, ctx.Path)
		http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
		return
	}

	// route matching
	entry, match := app.Find(ctx.Path)
	if entry == nil {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	handle := reflect.New(entry)
	exec, ok := handle.Interface().(Handler)
	if !ok {
		panic("exec is not Handler")
	}
	exec.init(app, ctx.ResponseWriter, ctx.Request)
	exec.parse_arguments(match)

	exec.Initialize()
	// check if status is not default value 0, Initialize is finished handler
	if ctx.Finished() {
		return
	}

	exec.Prepare()
	// check if status is not default value 0, Prepare is finished handler
	if ctx.Finished() {
		return
	}

	method := handle.MethodByName(ctx.Method)
	if (method == reflect.Value{}) {
		ctx.HTTPError(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	method.Call(nil)
	exec.Finish()
	if ctx.Finished() {
		return
	}
	// Finish handler request normally, set statusOK
	exec.WriteHeader(http.StatusOK)
}
