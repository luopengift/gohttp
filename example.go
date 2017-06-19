package gohttp

import (
	"fmt"
)

type TplHandler struct {
	HttpHandler
}

func (self *TplHandler) GET() {
	//self.Redirect("http://www.baidu.com", 301)
	self.Render("template/index.tpl", map[string]string{"content": "This is a test page"})
}

func (self *TplHandler) POST() {
	fmt.Println(self.Header())
	self.Output(self.GetBodyArgs())
}

type ArgsHandler struct {
	HttpHandler
}

func (self *ArgsHandler) Run() {
	match := self.GetMatchArgs()
	query := self.GetQueryArgs()
	body := self.GetBodyArgs()
	result := map[string]interface{}{
		"match":  match,
		"query":  query,
		"body":   string(body),
	}
	fmt.Println(result)
	err := self.Output(result)
    fmt.Println(err)
}

func (self *ArgsHandler) GET() {
	self.Run()
}

func (self *ArgsHandler) POST() {
	self.Run()
}

func init() {
	RouterRegister("^/args(/(?P<args>[0-9a-zA-Z]*))?$", &ArgsHandler{})
	RouterRegister("^/tpl", &TplHandler{})
}
