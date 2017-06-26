package gohttp

import (
	"github.com/luopengift/golibs/logger"
	"html/template"
	"io/ioutil"
	"net/http"
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

	// native http.ResponseWrite
	http.ResponseWriter

	// response is end or not
	isEnd bool

	// response header map
	Header map[string]string
	// response status code
	status int

	// contain datas need response to client
	response []byte
}

func NewHttpHandler(app *Application, responsewriter http.ResponseWriter, request *http.Request) *HttpHandler {
	httphandler := new(HttpHandler)
	httphandler.init(app, responsewriter, request)
	return httphandler
}

func (ctx *HttpHandler) init(app *Application, responsewriter http.ResponseWriter, request *http.Request) {
	ctx.Application = app
	ctx.Request = request
	ctx.Method = request.Method
	ctx.URL = request.RequestURI
	ctx.Remote = strings.Split(request.RemoteAddr, ":")[0]
	ctx.Path = request.URL.Path
	ctx.match = make(map[string]string)
	ctx.query = make(map[string][]string)
	ctx.body = []byte{}
	ctx.form = make(map[string][]string)
	ctx.ResponseWriter = responsewriter
	ctx.isEnd = false
	ctx.Header = make(map[string]string)
	ctx.status = http.StatusOK
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
	if ctx.isEnd {
		logger.Error("HttpHandler is end!")
		return
	}
	ctx.status = code
	http.Redirect(ctx.ResponseWriter, ctx.Request, url, code)
	ctx.isEnd = true
}

// response Http Error
func (ctx *HttpHandler) HTTPError(error string, code int) {
	if ctx.isEnd {
		logger.Error("HttpHandler is end!")
		return
	}
	ctx.status = code
	http.Error(ctx.ResponseWriter, error, code)
	ctx.isEnd = true
}

// If response is sent, do not sent again
func (ctx *HttpHandler) Output(v interface{}) {
	if ctx.isEnd {
		logger.Error("HttpHandler is end!")
		return
	}
	if response, err := ToBytes(v); err != nil {
		panic(err)
	} else {
		ctx.output(response)
		ctx.isEnd = true
	}
}

// response data from server to client
func (ctx *HttpHandler) output(response []byte) {
	for name, value := range ctx.Header {
		ctx.ResponseWriter.Header().Set(name, value)
	}
	ctx.ResponseWriter.WriteHeader(ctx.status)
	ctx.ResponseWriter.Write(response)
}

// If response is sent, do not sent again
func (ctx *HttpHandler) Render(tpl string, data interface{}) {
	if ctx.isEnd {
		logger.Error("HttpHandler is end!")
		return
	}
	if template, ok := (*ctx.Template)[tpl]; !ok {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		ctx.render(template, data)
		ctx.isEnd = true
	}
}

// render html data to client
func (ctx *HttpHandler) render(tpl *template.Template, data interface{}) {
	for name, value := range ctx.Header {
		ctx.ResponseWriter.Header().Set(name, value)
	}
	ctx.ResponseWriter.WriteHeader(ctx.status)
	(*tpl).Execute(ctx.ResponseWriter, data)
}

// set response header into Header
func (ctx *HttpHandler) SetHeader(name, value string) {
	ctx.Header[name] = value
}

// set response status code int status
func (ctx *HttpHandler) SetStatusCode(code int) {
	ctx.status = code
}

func (ctx *HttpHandler) Prepare() {}
func (ctx *HttpHandler) Finish()  {}
