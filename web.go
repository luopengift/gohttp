package gohttp

import (
    "io/ioutil"
    "fmt"
	"net/http"
	"reflect"
	"time"
    "sync"
)

var (
	SUPPORTED_METHODS = [7]string{
		"GET",
		"HEAD",
		"POST",
		"DELETE",
		"PATCH",
		"PUT",
		"OPTIONS",
	}
)

type RequestHandler struct {
    //query_arguments []string
    //get_arguments   []string
    //post_arguments  []string
    matchArgs   map[string]string
    queryArgs map[string][]string
    bodyArgs    []byte
    finished bool
}

func (self *HttpHandler) Header() http.Header {
    return self.Request().Header
}

func (self *HttpHandler) GetHeader(name string) string {
    if value,ok := self.Request().Header[name];ok {
        return value[0]
    }
    return ""
}

func (self *HttpHandler) GetMatchArgs() map[string]string {
    return self.matchArgs
}

func (self *HttpHandler) GetMatchArg(name string) string {
    if value,ok := self.matchArgs[name];ok {
        return value
    }
    return ""
}

func (self *HttpHandler) GetQueryArgs() map[string][]string {
    return self.queryArgs
}

func (self *HttpHandler) GetQueryArg(name string) string {
    if value,ok := self.queryArgs[name];ok {
        return value[0]
    }
    return ""
}
//Content-Type:text/plain;charset=UTF-8
func (self *HttpHandler) GetBodyArgs() []byte {
    return self.bodyArgs
}

func (self *HttpHandler) GetBodyArg(name string) {

}


//Content-Type:"application/x-www-form-urlencoded"

func (self *HttpHandler) GetFormArgs() {}
func (self *HttpHandler) GetFormArg() {}


type Handler interface {
	http.Handler
    Connection
    Init(*Conn,map[string]string)
}

type HttpHandler struct {
    sync.Pool
    RequestHandler
	*Conn
}

func NewHttpHandler() *HttpHandler {
    httphandler := &HttpHandler{}
    httphandler.Pool.New = func() interface{} {
        return &Conn{}
    }
    return httphandler
}


func (self *HttpHandler) Init(conn *Conn,kv map[string]string) {
	self.Conn = conn
    self.matchArgs = kv  //获取通过正则匹配出来的url参数
    self.Request().ParseForm()
    self.queryArgs = self.Request().Form  //获取query参数

    self.bodyArgs, _ = ioutil.ReadAll(self.Request().Body) //获取body参数
}

func (self *HttpHandler) Output(o []byte) {
    self.Response().Write(o)
}



func (self *HttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	stime := time.Now()

    if match,entry := self.findHandle(request.URL.Path);match == nil {
		http.Error(response, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}else{
	    conn := self.Pool.Get().(*Conn)
	    defer self.Pool.Put(conn)
	    conn.init(response, request)

        handle := reflect.New(entry)
        handle.Interface().(Handler).Init(conn,match)
        handle.MethodByName("Prepare").Call(nil)
        handle.MethodByName(request.Method).Call(nil)
        handle.MethodByName("Finish").Call(nil)
    }
    fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), 200,request.Method, request.URL, time.Since(stime))
}

func (self *HttpHandler) findHandle(url string) (map[string]string,muxEntry) {
    for pattern,handle := range Router {
        if match := pattern.FindStringSubmatch(url);match != nil {
            var kv = map[string]string{}
            for key,value := range pattern.SubexpNames() {
                kv[value] = match[key]
            }
            return kv,handle
        }
    }
    return nil,nil
}
