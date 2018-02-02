package gohttp

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/types"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// Handler implements http handler interface.
// Prepare -> Auth -> GET/POST... -> Finish
type Handler interface {
	http.Handler
	// Prepare invoked before Init.
	Prepare()
	// Auth invoked before httpMethod func.
	Auth()
	// Finish invoked after httpMethod func.
	Finish()
	parse_arguments(match map[string]string)
	init(*Application, http.ResponseWriter, *http.Request)
}

type HttpHandler struct {
	// application
	*Application

	// request based native *http.Request
	*request

	// ResponseEriter based native http.ResponseWrite
	// Implements http.ResponseWriter interface and some extra interface,
	// Such as, Status() int, Finished() bool, Size() int
	ResponseWriter
}

func NewHttpHandler(app *Application, responsewriter http.ResponseWriter, request *http.Request) *HttpHandler {
	httphandler := new(HttpHandler)
	httphandler.init(app, responsewriter, request)
	return httphandler
}

func (ctx *HttpHandler) init(app *Application, responsewriter http.ResponseWriter, request *http.Request) {
	ctx.Application = app
	ctx.request = NewRequestReader(request)
	ctx.ResponseWriter = NewResponseWriter(responsewriter)
}

// App returns *Application instance in this HttpHandler context.
func (ctx *HttpHandler) App() *Application {
	return ctx.Application
}

func cookie(name, value string, expire int) *http.Cookie {
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		MaxAge:  expire,
		Expires: time.Now().Add(time.Duration(expire) * time.Second),
	}
	return cookie
}

// set cookie for response
func (ctx *HttpHandler) SetCookie(name, value string) {
	cookie := cookie(name, value, 86400)
	http.SetCookie(ctx.ResponseWriter, cookie)
}

func (ctx *HttpHandler) GetQueryArgs() map[string][]string {
	return ctx.query
}

// fetch query argument named by <name>, null is default value defined by user
func (ctx *HttpHandler) GetQuery(name string, null string) string {
	if value, ok := ctx.query[name]; !ok {
		return null
	} else {
		return value[0]
	}
}

func (ctx *HttpHandler) GetMatchArgs() map[string]string {
	return ctx.match
}

// fetch match argument named by <name>, null is default value defined by user
func (ctx *HttpHandler) GetMatch(name string, null string) string {
	if value, ok := ctx.match[name]; !ok {
		return null
	} else {
		return value
	}
}

// fetch body arguments
func (ctx *HttpHandler) GetBodyArgs() []byte {
	return ctx.body
}

// fetch body argument named by <name>
func (ctx *HttpHandler) GetBody(name string) interface{} {
	if body, err := types.BytesToMap(ctx.body); err != nil {
		panic(err)
	} else {
		return body[name]
	}
}

// fetch form argument named by <name>
func (ctx *HttpHandler) GetForm(name string) string {
	//TODO
	return ""
}

// response redirect
func (ctx *HttpHandler) Redirect(url string, code int) {
	ctx.WriteHeader(code)
	http.Redirect(ctx.ResponseWriter, ctx.Request, url, code)
}

// response Http Error
func (ctx *HttpHandler) HTTPError(error string, code int) {

	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.ResponseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.output([]byte(error), code)

}

// Output response the http request. If response is sent, do not sent again.
func (ctx *HttpHandler) Output(v interface{}, code ...int) {
	if ctx.Finished() {

		return
	}
	response, err := types.ToBytes(v)
	if err != nil {
		panic(err)
	}
	if len(code) == 0 {
		ctx.output(response, 200)
	} else {
		ctx.output(response, code[0])
	}
}

// response data from server to client
func (ctx *HttpHandler) output(response []byte, code int) {
	//for name, value := range ctx.Header {
	//	ctx.ResponseWriter.Header().Set(name, value)
	//}
	ctx.WriteHeader(code)
	ctx.ResponseWriter.Write(response)
}

// If response is sent, do not sent again
func (ctx *HttpHandler) Render(tpl string, data interface{}) {
	path := filepath.Join(ctx.Config.WebPath, tpl)

	//TODO:check it twice,not a good choice
	if _, ok := (*ctx.Template)[path]; !ok {
		(*ctx.Template).AddFile(path)
	}
	if template, ok := (*ctx.Template)[path]; !ok {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		ctx.render(template, data)
		return
	}
}

// render html data to client
func (ctx *HttpHandler) render(tpl *template.Template, data interface{}) {
	//for name, value := range ctx.Header {
	//	ctx.ResponseWriter.Header().Set(name, value)
	//}
	ctx.WriteHeader(http.StatusOK) //200
	(*tpl).Execute(ctx.ResponseWriter, data)
}

// file download response by file path.
func (ctx *HttpHandler) Download(file string) {
	if ctx.Finished() {
		logger.Error("HttpHandler is end!")
		return
	}
	f, err := os.Stat(file)
	if err != nil {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound) //404
		return
	}
	if f.IsDir() {
		ctx.HTTPError(http.StatusText(http.StatusForbidden), http.StatusForbidden) //403
		return
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+path.Base(file))
	ctx.ResponseWriter.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.ResponseWriter.Header().Set("Expires", "0")
	ctx.ResponseWriter.Header().Set("Cache-Control", "must-revalidate")
	ctx.ResponseWriter.Header().Set("Pragma", "public")
	http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
	return
}

func (ctx *HttpHandler) Prepare() {}
func (ctx *HttpHandler) Auth()    {}
func (ctx *HttpHandler) Finish()  {}
