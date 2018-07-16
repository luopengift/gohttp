package gohttp

import (
	"net/http"
	"time"
)

// Handler implements http handler interface.
// Initialize -> Prepare  -> GET/POST... -> Finish
type Handler interface {
	// Handler implements ServeHTTP(http.ResponseWriter, *http.Request) func.
	http.Handler
	// Prepare invoked before Init.
	Prepare()
	// Initialize invoked before httpMethod func.
	Initialize()
	// Finish invoked after httpMethod func.
	Finish()
	WriteHeader(code int)
	parseArgs() error
	init(*Context)
}

// BaseHTTPHandler http handler
type BaseHTTPHandler struct {
	*Context
}

func (ctx *BaseHTTPHandler) init(context *Context) {
	ctx.Context = context
}

func cookie(name, value string, expire int) *http.Cookie {
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		MaxAge:  expire,
		Expires: time.Now().Add(time.Duration(expire) * time.Second),
	}
	return cookie
}

// SetCookie set cookie for response
func (ctx *BaseHTTPHandler) SetCookie(name, value string) {
	cookie := cookie(name, value, 86400)
	http.SetCookie(ctx.ResponseWriter, cookie)
}

// Initialize init
func (ctx *BaseHTTPHandler) Initialize() {}

// Prepare xx
func (ctx *BaseHTTPHandler) Prepare() {}

// HEAD method
func (ctx *BaseHTTPHandler) HEAD() {}

// Finish func
func (ctx *BaseHTTPHandler) Finish() {}
