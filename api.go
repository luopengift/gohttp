package gohttp

// ApiOutput is sturct data need responsed.
type ApiOutput struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Err  error       `json:"err"`
	Data interface{} `json:"data"`
}

// ApiHandler designed for http api. It can used easy.
type ApiHandler struct {
	ApiOutput
	HttpHandler
}

func (ctx *ApiHandler) Finish() {
	ctx.Output(ctx.ApiOutput)
}
