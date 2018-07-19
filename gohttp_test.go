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
	app.RouteFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo ok"))
	})
	app.RouteFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://www.baidu.com", 301)
	})
	app.RouteFunCtx("/bar", func(ctx *Context) {
		ctx.Output("bar ok")
	})
	app.RouteAlias("/fooalias", "/foo")
	app.RouteAlias("/ttt", "/fooalias")
	app.Route("/baz", &baz{})
	app.SetTLS("cert.pem", "key.pem")
	app.Info("%#v", app.Config)
	app.Run(":8888")
}
