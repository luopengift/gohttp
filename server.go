package gohttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Handler interface {
	http.Handler
	Connection
	Init(*Conn, map[string]string)
}

type HttpHandler struct {
	status int64
	match  map[string]string
	query  map[string][]string
	body   []byte

	sync.Pool
	RequestHandler
}

func (self *HttpHandler) Header() http.Header {
	return self.Request.Header
}

func (self *HttpHandler) GetHeader(name string, null ...string) string {
	if value, ok := self.Request.Header[name]; ok {
		return value[0]
	}
	//默认值
	if len(null) == 1 {
		return null[0]
	}
	return ""
}

func (self *HttpHandler) GetMatchArgs() map[string]string {
	return self.match
}

//获取通过正则表达式匹配到的uri中的参数
//name:参数名, null:默认值
func (self *HttpHandler) GetMatchArg(name string, null ...string) string {
	if value, ok := self.match[name]; ok {
		return value
	}
	//默认值
	if len(null) == 1 {
		return null[0]
	}

	return ""
}

func (self *HttpHandler) GetQueryArgs() map[string][]string {
	return self.query
}

func (self *HttpHandler) GetQueryArg(name string, null ...string) string {
	if value, ok := self.query[name]; ok {
		return value[0]
	}
	//默认值
	if len(null) == 1 {
		return null[0]
	}
	return ""
}

//Content-Type:text/plain;charset=UTF-8
func (self *HttpHandler) GetBodyArgs() []byte {
	return self.body
}

func (self *HttpHandler) GetBodyArg(name string, null ...string) interface{} {
	if body, err := BytesToJson(self.body); err != nil {
	} else {
		return body[name]
	}

	//默认值
	if len(null) == 1 {
		return null[0]
	}
	return ""
}

//Content-Type:"application/x-www-form-urlencoded"

func (self *HttpHandler) GetFormArgs() {}
func (self *HttpHandler) GetFormArg()  {}

func NewHttpHandler() *HttpHandler {
	httphandler := &HttpHandler{}
	httphandler.Pool.New = func() interface{} {
		return &Conn{}
	}
	return httphandler
}

func (self *HttpHandler) Init(conn *Conn, kv map[string]string) {
	self.status = http.StatusOK
	self.Conn = conn
	self.match = kv //获取通过正则匹配出来的url参数
	self.Request.ParseForm()
	self.query = self.Request.Form //获取query参数

	self.body, _ = ioutil.ReadAll(self.Request.Body) //获取body参数
}

func (self *HttpHandler) Output(o interface{}) error {
	out, err := ToBytes(o)
	self.ResponseWriter.Write(out)
	return err
}

func (self *HttpHandler) ServeHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	stime := time.Now()
	self.status = http.StatusOK
	if strings.HasPrefix(request.URL.Path, "/static") || hasSuffixs(request.URL.Path, ".ico", ".jpg", ".jpeg", ".png", ".bmp", ".gif", ".js", ".css", ".swf") {
		StaticPath := "."
		file := filepath.Join(StaticPath, request.URL.Path)
		http.ServeFile(responsewriter, request, file)
		goto END
	}

	if match, entry := findHandle(request.URL.Path); match == nil {
		self.status = http.StatusNotFound
		http.Error(responsewriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		goto END
	} else {
		conn := self.Pool.Get().(*Conn)
		defer self.Pool.Put(conn)
		conn.init(responsewriter, request)

		handle := reflect.New(entry)
		handle.Interface().(Handler).Init(conn, match)
		handle.MethodByName(request.Method).Call(nil)
		self.status = reflect.Indirect(handle).FieldByName("status").Int()
		goto END
	}
END:
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), self.status, request.Method, request.URL, request.RemoteAddr, "->", request.Host, time.Since(stime))
}

func findHandle(url string) (map[string]string, muxEntry) {
	for pattern, handle := range RouterMap {
		if match := pattern.FindStringSubmatch(url); match != nil {
			var kv = map[string]string{}
			for key, value := range pattern.SubexpNames() {
				kv[value] = match[key]
			}
			delete(kv, "")
			return kv, handle
		}
	}
	return nil, nil
}

func (self *HttpHandler) Render(tpl string, data ...interface{}) (err error) {
	if len(data) == 1 {
		self.status, err = renderFile(self.ResponseWriter, tpl, data[0])
		return
	}
	self.status, err = renderFile(self.ResponseWriter, tpl, nil)
	return
}

func (self *HttpHandler) ReanderString(name string, data interface{}) error {
	return renderString(self.ResponseWriter, name, data)

}
