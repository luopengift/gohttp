package gohttp

import (
	"fmt"
	"strings"
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
	result := []string{fmt.Sprintf("\nMethod: %s", ctx.Request.Method),
		fmt.Sprintf("Protocol: %s", ctx.Request.Proto),
		fmt.Sprintf("Host: %s", ctx.Request.Host),
		fmt.Sprintf("RemoteAddr: %s", ctx.Request.RemoteAddr),
		fmt.Sprintf("RequestURI: %q", ctx.Request.RequestURI),
		fmt.Sprintf("URL: %#v", ctx.Request.URL),
		fmt.Sprintf("Body.ContentLength: %d (-1 means unknown)", ctx.Request.ContentLength),
		fmt.Sprintf("Close: %v (relevant for HTTP/1 only)", ctx.Request.Close),
		fmt.Sprintf("TLS: %#v", ctx.Request.TLS),
		fmt.Sprintf("\nHeaders: \n"),
	}
	ctx.Output(strings.Join(result, "\n"))
	ctx.Request.Header.Write(ctx.ResponseWriter)
}
