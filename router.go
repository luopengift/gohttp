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
	method string
	path   string
	alias  string
	regx   *regexp.Regexp
	entry  HandleHTTP
}

func (r route) String() string {
	return fmt.Sprintf("method:%s, path:%s, alias:%s, entry: %v", r.method, r.path, r.alias, r.entry)
}

// RouterList router List
type RouterList []*route

// InitRouterList init route list
func InitRouterList() *RouterList {
	return new(RouterList)
}

func (r *RouterList) append(method, path, alias string, entry HandleHTTP) {
	route := &route{method: method, path: path, alias: alias, regx: regexp.MustCompile(path), entry: entry}
	*r = append(*r, route)
}

// Route route
func (r *RouterList) Route(path string, handler Handler) {
	rv := reflect.ValueOf(handler)
	rt := reflect.Indirect(rv).Type()
	entry := Entry{rt}
	r.append("", path, "", entry)
}

// RouteFunc route handle func
func (r *RouterList) RouteFunc(path string, f HandleFunc) {
	r.append("", path, "", f)
}

// RouteFunCtx route handle func
func (r *RouterList) RouteFunCtx(path string, f HandleFunCtx) {
	r.append("", path, "", f)
}

// RouteMethod route by method
func (r *RouterList) RouteMethod(method, path string, f HandleFunc) {
	r.append(method, path, "", f)
}

// RouteCtxMethod route by method
func (r *RouterList) RouteCtxMethod(method, path string, f HandleFunCtx) {
	r.append(method, path, "", f)
}

// RouteAlias alias path
func (r *RouterList) RouteAlias(path, alias string) {
	r.append("", path, alias, nil)
}

// find search route
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
