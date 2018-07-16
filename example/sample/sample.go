package main

import (
	"net/http"

	"github.com/luopengift/gohttp"
)

type baz struct {
	gohttp.BaseHTTPHandler
}

func (ctx *baz) GET() {
	ctx.Output("baz ok")
}

func main() {
	app := gohttp.Init()
	// register route "/foo"
	app.RouteFunc("/foo", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("foo ok"))
	})
	// register route "/bar"
	app.RouteFunCtx("/bar", func(ctx *gohttp.Context) {
		ctx.Output("bar ok")
	})
	// register route "/baz"
	app.Route("/baz", &baz{})
	app.Run(":8888")
}
