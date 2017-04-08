package gohttp

import (
    "fmt"
)

type Test struct {
	HttpHandler
}

func (self *Test) GET() {
	//self.Redirect("http://www.baidu.com", 301)
	self.Render("template/index.tpl", map[string]string{"content":"This is a test page"})
}

func (self *Test) POST() {
	self.Output([]byte(self.Request.PostFormValue("id")))
	self.Output([]byte("hello"))
}

type RouterHandler struct {
    HttpHandler
}

func (self *RouterHandler) GET() {
    for k,v := range RouterMap {
        self.Output([]byte(fmt.Sprintf("%v:%v\n",k,v)))
    }
}

func init() {
	RouterRegister("^/routers$",&RouterHandler{})
	RouterRegister("^/(?P<ID>[0-9]*)/(?P<NAME>[a-zA-Z]*)$", &Test{})
	RouterRegister("^/test", &Test{})
}
