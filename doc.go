// Package gohttp is sample http server framework
// This package is
// 1. used like std net/http
// app.Init()
// app.RouteFunc("/foo", func(resp http.ResponseWriter, req http.Request) {
// 	   resp.Write([]byte("foo ok")
// })
//
// 2. used gohttp.Context
// app.Init()
// app.RouteFunCtx("/bar", func(ctx *gohttp.Context) {
//     ctx.Output("bar ok")
// })
// 3. used like tornado
// type baz struct {
//     gohttp.Context
// }
// func (ctx *baz) GET() {
//  ctx.Output("baz ok")
// }
// app.Init()
// app.Route("/baz", &baz{})
//
package gohttp
