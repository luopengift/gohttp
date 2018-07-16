package main

import (
	"fmt"
	"net/http"

	"github.com/luopengift/gohttp"
)

// RoutersHandler routers handler
type RoutersHandler struct {
	gohttp.BaseHTTPHandler
}

// GET method
func (ctx *RoutersHandler) GET() {
	ctx.Output(ctx.RouterList.String())
	//ctx.Redirect("http://www.baidu.com", 301)
	//ctx.Render("template/index.tpl", map[string]string{"content": "This is a test page"})
}

// POST method
func (ctx *RoutersHandler) POST() {
	fmt.Println("tpl post")
	ctx.Output("tpl post")
}

// MirrorHandler mirror handler
type MirrorHandler struct {
	gohttp.BaseHTTPHandler
}

// Prepare func
func (ctx *MirrorHandler) Prepare() {
	if ctx.Method == "GET" {
		ctx.Output("hello Prepare inject", 401)
	}
	if ctx.Method == "PUT" {
		panic(http.ErrAbortHandler)
	}
}

// Run handler entry
func (ctx *MirrorHandler) Run() {
	match := ctx.GetMatchArgs()
	query := ctx.GetQueryArgs()
	body := ctx.GetBodyArgs()
	result := map[string]interface{}{
		"match": match,
		"query": query,
		"body":  string(body),
	}
	fmt.Println(result)
	ctx.Output(result)
}

// GET method
func (ctx *MirrorHandler) GET() {
	ctx.Run()
}

// POST method
func (ctx *MirrorHandler) POST() {
	ctx.Run()
}

func main() {
	app := gohttp.Init()
	app.Route("^/mirror(/(?P<args>[0-9a-zA-Z]*))?$", &MirrorHandler{})
	app.Route("^/routers$", &RoutersHandler{})
	app.Run(":18081")
}
