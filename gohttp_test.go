package gohttp

import (
	"net/http"
	"testing"
)

type baz struct {
	BaseHTTPHandler
}

func (ctx *baz) GET() {
	ctx.Output("baz ok")
}

func Test_http(t *testing.T) {
	app := Init()
	app.RouteFunc("/foo", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("foo ok"))
	})
	app.RouteFunCtx("/bar", func(ctx *Context) {
		ctx.Output("bar ok")
	})
	app.Route("/baz", &baz{})
	app.SetTLS("cert.pem", "key.pem")
	app.Run(":8888")
}
