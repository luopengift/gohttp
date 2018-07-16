package gohttp

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
)

// HandleHTTP handle http interface
type HandleHTTP interface {
	Exec(*Context)
}

// HandleFunc handle func
type HandleFunc func(http.ResponseWriter, *http.Request)

// Exec implements HandleHTTP interface
func (f HandleFunc) Exec(ctx *Context) {
	f(ctx.ResponseWriter, ctx.Request)
}

// HandleFunCtx handle fun ctx
type HandleFunCtx func(*Context)

// Exec implements HandleHTTP interface
func (f HandleFunCtx) Exec(ctx *Context) {
	f(ctx)
}

// Entry xx
type Entry struct {
	reflect.Type
}

// Exec implements HandleHTTP interface
func (entry Entry) Exec(ctx *Context) {
	handle := reflect.New(entry.Type)
	exec, ok := handle.Interface().(Handler)
	if !ok {
		panic("exec is not Handler")
	}
	exec.init(ctx)
	if err := exec.parseArgs(); err != nil {
		ctx.Warn("parse args error: %v", err)
		return
	}

	exec.Initialize()
	// check if status is not default value 0, Initialize is finished handler
	if ctx.Finished() {
		return
	}

	exec.Prepare()
	// check if status is not default value 0, Prepare is finished handler
	if ctx.Finished() {
		return
	}

	method := handle.MethodByName(ctx.Method)
	if (method == reflect.Value{}) {
		ctx.HTTPError(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	method.Call(nil)
	exec.Finish()
	if ctx.Finished() {
		return
	}
	// Finish handler request normally, set statusOK
	exec.WriteHeader(http.StatusOK)
}

type muxEntry reflect.Type

type route struct {
	path   string
	method string
	regx   *regexp.Regexp
	entry  HandleHTTP
}

func newroute(path string, handler Handler) *route {
	rv := reflect.ValueOf(handler)
	rt := reflect.Indirect(rv).Type()
	return &route{path: path, regx: regexp.MustCompile(path), entry: Entry{rt}}
}

// RouterList router List
type RouterList []*route

// InitRouterList init route list
func InitRouterList() *RouterList {
	return new(RouterList)
}

// Route route
func (r *RouterList) Route(path string, handler Handler) {
	route := newroute(path, handler)
	*r = append(*r, route)
}

// RouteFunc route handle func
func (r *RouterList) RouteFunc(path string, f HandleFunc) {
	route := &route{path: path, regx: regexp.MustCompile(path), entry: f}
	*r = append(*r, route)

}

// RouteFunCtx route handle func
func (r *RouterList) RouteFunCtx(path string, f HandleFunCtx) {
	route := &route{path: path, regx: regexp.MustCompile(path), entry: f}
	*r = append(*r, route)
}

func (r *RouterList) find(path string) (*route, map[string]string) {
	for _, route := range *r {
		if match := route.regx.FindStringSubmatch(path); match != nil {
			kv := make(map[string]string)
			for key, value := range route.regx.SubexpNames() {
				kv[value] = match[key]
			}
			delete(kv, "")
			return route, kv
		}
	}
	return nil, nil
}

func (r *RouterList) String() string {
	str := "\nRouter Map:\n"
	for _, route := range *r {
		str += fmt.Sprintf("%v => %v\n", route.path, route.entry)
	}
	return str
}
