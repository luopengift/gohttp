package gohttp

import (
	"github.com/luopengift/golibs/logger"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Handler interface {
	http.Handler
	Prepare()
	Finish()
	parse_arguments(match map[string]string)
	init(*Application, http.ResponseWriter, *http.Request)
}

/*--------------*/

type HttpHandler struct {
	// application
	*Application
	// native *http.Request
	*http.Request

	// ResponseEriter based native http.ResponseWrite
	// Implements http.ResponseWriter interface and some extra interface,
	// Such as, Status() int, Finished() bool, Size() int
	ResponseWriter

	// request method
	Method string

	// request url
	URL string

	// request host without port
	Remote string
	//
	Path string
	// request path regx match arguments
	match map[string]string
	//request query arguments
	query map[string][]string
	// request body arguments
	body []byte
	// TODO:request form arguments
	form map[string][]string
}

func NewHttpHandler(app *Application, responsewriter http.ResponseWriter, request *http.Request) *HttpHandler {
	httphandler := new(HttpHandler)
	httphandler.init(app, responsewriter, request)
	return httphandler
}

func (ctx *HttpHandler) init(app *Application, responsewriter http.ResponseWriter, request *http.Request) {
	ctx.Application = app
	ctx.Request = request
	ctx.ResponseWriter = NewResponseWriter(responsewriter)
	ctx.Method = request.Method
	ctx.URL = request.RequestURI
	ctx.Remote = strings.Split(request.RemoteAddr, ":")[0]
	ctx.Path = request.URL.Path
	ctx.match = make(map[string]string)
	ctx.query = make(map[string][]string)
	ctx.body = []byte{}
	ctx.form = make(map[string][]string)
}

// App returns *Application instance in this HttpHandler context.
func (ctx *HttpHandler) App() *Application {
	return ctx.Application
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
	if body, err := BytesToJson(ctx.body); err != nil {
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

// prepare match and assignment to match arguments
func (ctx *HttpHandler) prepare_match_arguments(match map[string]string) {
	ctx.match = match
}

// prepare query and assignment to query arguments
func (ctx *HttpHandler) prepare_query_arguments() {
	ctx.query = ctx.Request.Form
}

// prepare body and assignment to body arguments
func (ctx *HttpHandler) prepare_body_arguments() {
	var err error
	ctx.body, err = ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}
}

// prepare form and assignment to form arguments
// Content-Type:application/x-www-form-urlencoded
func (ctx *HttpHandler) prepare_form_arguments() {
	ctx.form = ctx.Request.PostForm
}

// parse and handler arguments
func (ctx *HttpHandler) parse_arguments(match map[string]string) {
	// parse form automatically
	ctx.Request.ParseForm()

	ctx.prepare_match_arguments(match)
	ctx.prepare_query_arguments()
	ctx.prepare_body_arguments()
	ctx.prepare_form_arguments()
	//logger.Debug("header:%#v", ctx.Request.Header)
	//logger.Debug("match:%#v,query:%#v,body:%#v", ctx.match, ctx.query, ctx.body)
	//logger.Debug("PostForm:%#v,MultipartForm:%#v", ctx.Request.PostForm, ctx.Request.MultipartForm)
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

// If response is sent, do not sent again
func (ctx *HttpHandler) Output(v interface{}) {
	if response, err := ToBytes(v); err != nil {
		panic(err)
	} else {
		ctx.output(response, 200)
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
	path := filepath.Join(ctx.Config.StaticPath, tpl)
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
func (ctx *HttpHandler) Finish()  {}
