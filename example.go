package gohttp

import ()

type Test struct {
	HttpHandler
}

func (self *Test) GET() {
	//	self.Redirect("http://www.baidu.com", 301)
	self.Render("template/index.tpl", nil)
}

func (self *Test) POST() {
	self.Output([]byte(self.Request.PostFormValue("id")))
	self.Output([]byte("hello"))
}

func init() {
	RouterRegister("^/(?P<ID>[0-9]*)/(?P<NAME>[a-zA-Z]*)$", &Test{})
}
