package gohttp

import (
	"fmt"
	"github.com/luopengift/golibs/logger"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

type Application struct {
	*Config
	*Template
	*RouterList
	*http.Server
}

func Init() *Application {
	app := new(Application)
	app.Config = InitConfig()
	app.Template = InitTemplate()
	app.RouterList = InitRouterList()
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
	return app
}

func (app *Application) Run(addr ...string) {
	if len(addr) != 0 {
		app.Server.Addr = addr[0]
	} else {
		app.Server.Addr = app.Config.Addr
	}
	fmt.Println("HttpsServer Start", app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// ServeHTTP is HTTP server implement method. It makes App compatible to native http handler.
func (app *Application) ServeHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	app.handler(responsewriter, request)
}

func (app *Application) handler(responsewriter http.ResponseWriter, request *http.Request) {
	stime := time.Now()

	// LOG Format
	format := "%3d %s %s (%s) %s"

	// init a new http handler
	ctx := NewHttpHandler(app, responsewriter, request)
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			ctx.HTTPError(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) //500
			logger.Error(format+" | %v", ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime), err)
		}
	}()

	// handler static file
	if strings.HasPrefix(ctx.Path, "/static") || hasSuffixs(ctx.Path) {
		file := filepath.Join(ctx.Config.StaticPath, ctx.Path)
		//TODO BUG: if file is not found, log http status is 200
		http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
		goto END
	}

	// route matching
	if entry, match := app.Find(ctx.Path); entry == nil {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound)
		goto END //404
	} else {
		handle := reflect.New(entry)
		exec, ok := handle.Interface().(Handler)
		if !ok {
			panic("exec is not Handler")
		}
		exec.init(app, ctx.ResponseWriter, ctx.Request)
		exec.parse_arguments(match)
		exec.Prepare()

		// check if status is not default value 0, knows prepare is finished handler
		if ctx.Finished() {
			goto END //Finished
		}
		if method := handle.MethodByName(ctx.Method); bool(method == reflect.Value{}) {
			ctx.HTTPError(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			goto END //405
		} else {
			method.Call(nil)
			exec.Finish()
			goto END
		}
	}
END:
	switch ctx.Status() {
	case 200, 301:
		logger.Info(format, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	case 400, 401, 403, 404, 405:
		logger.Warn(format, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	case 500, 501, 502, 503:
		logger.Error(format, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	default:
		logger.Error(format, ctx.Status(), ctx.Method, ctx.URL, ctx.Remote, time.Since(stime))
	}
}
