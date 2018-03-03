package gohttp

import (
	"fmt"
)

type RouteHandler struct {
	HttpHandler
}

func (ctx *RouteHandler) GET() {
	ctx.Output(ctx.RouterList.String())
}

type InfoHandler struct {
	HttpHandler
}

func (ctx *InfoHandler) GET() {
	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
	result := fmt.Sprintf("Method: %s\n", ctx.Request.Method)
	result += fmt.Sprintf("Protocol: %s\n", ctx.Request.Proto)
	result += fmt.Sprintf("Host: %s\n", ctx.Request.Host)
	result += fmt.Sprintf("RemoteAddr: %s\n", ctx.Request.RemoteAddr)
	result += fmt.Sprintf("RequestURI: %q\n", ctx.Request.RequestURI)
	result += fmt.Sprintf("URL: %#v\n", ctx.Request.URL)
	result += fmt.Sprintf("Body.ContentLength: %d (-1 means unknown)\n", ctx.Request.ContentLength)
	result += fmt.Sprintf("Close: %v (relevant for HTTP/1 only)\n", ctx.Request.Close)
	result += fmt.Sprintf("TLS: %#v\n", ctx.Request.TLS)
	result += fmt.Sprintf("\nHeaders: \n")
	ctx.Output(result)
	ctx.Request.Header.Write(ctx.ResponseWriter)
}
