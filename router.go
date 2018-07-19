package gohttp

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
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

// Entry handle implements HandleHTTP interface
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

	if ctx.Finished() {
		return
	}
	exec.Finish()

}

type route struct {
	path   string
	method string
	alias  string
	regx   *regexp.Regexp
	entry  HandleHTTP
}

// RouterList router List
type RouterList []*route

// InitRouterList init route list
func InitRouterList() *RouterList {
	return new(RouterList)
}

// Route route
func (r *RouterList) Route(path string, handler Handler) {
	rv := reflect.ValueOf(handler)
	rt := reflect.Indirect(rv).Type()
	route := &route{path: path, regx: regexp.MustCompile(path), entry: Entry{rt}}
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

// RouteAlias alias path
func (r *RouterList) RouteAlias(path, targetPath string) {
	route := &route{path: path, regx: regexp.MustCompile(path), alias: targetPath}
	*r = append(*r, route)
}
func (r *RouterList) find(path string) (*route, map[string]string) {
	for _, route := range *r {
		if match := route.regx.FindStringSubmatch(path); match != nil {
			if route.alias != "" {
				return r.find(route.alias)
			}
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
	strList := make([]string, len(*r)+1)
	strList[0] = "\nRouter Map:"
	for idx, route := range *r {
		if route.alias != "" {
			strList[idx+1] = fmt.Sprintf("%v => %v\n", route.path, route.alias)
		} else {
			strList[idx+1] = fmt.Sprintf("%v : %v\n", route.path, route.entry)
		}
	}
	return strings.Join(strList, "")
}
