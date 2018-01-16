package gohttp

type RouteHandler struct {
	HttpHandler
}

func (ctx *RouteHandler) GET() {
	ctx.Output(ctx.RouterList.String())
}
