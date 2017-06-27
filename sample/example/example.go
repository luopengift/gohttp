package example

import (
	"fmt"
	"github.com/luopengift/gohttp"
	"net/http"
)

type RoutersHandler struct {
	gohttp.HttpHandler
}

func (ctx *RoutersHandler) GET() {
	ctx.Output(ctx.RouterList.String())
	//ctx.Redirect("http://www.baidu.com", 301)
	//ctx.Render("template/index.tpl", map[string]string{"content": "This is a test page"})
}

func (ctx *RoutersHandler) POST() {
	fmt.Println("tpl post")
	ctx.Output("tpl post")
}

type MirrorHandler struct {
	gohttp.HttpHandler
}

func (ctx *MirrorHandler) Prepare() {
	if ctx.Method == "GET" {
		ctx.Output("hello Prepare inject")
	}
	if ctx.Method == "PUT" {
		panic(http.ErrAbortHandler)
	}
}

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

func (ctx *MirrorHandler) GET() {
	ctx.Run()
}

func (ctx *MirrorHandler) POST() {
	ctx.Run()
}

var App *gohttp.Application

func init() {
	App = gohttp.Init()
	App.Route("^/mirror(/(?P<args>[0-9a-zA-Z]*))?$", &MirrorHandler{})
	App.Route("^/routers$", &RoutersHandler{})
}
