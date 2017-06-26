package gohttp

import (
	"fmt"
	"net/http"
)

type TplHandler struct {
	HttpHandler
}

func (self *TplHandler) GET() {
	self.Redirect("http://www.baidu.com", 301)
	//self.Render("template/index.tpl", map[string]string{"content": "This is a test page"})
	//fmt.Println("tpl get")
	//self.Output("tpl get")
}

func (self *TplHandler) POST() {
	fmt.Println("tpl post")
	self.Output("tpl post")
}

type MirrorHandler struct {
	HttpHandler
}

func (self *MirrorHandler) Prepare() {
	if self.Method == "GET" {
		self.Output("hello Prepare inject")
	}
	if self.Method == "PUT" {
		panic(http.ErrAbortHandler)
	}
}

func (self *MirrorHandler) Run() {
			match := self.GetMatchArgs()
			query := self.GetQueryArgs()
			body := self.GetBodyArgs()
			result := map[string]interface{}{
				"match": match,
				"query": query,
				"body":  string(body),
			}
			fmt.Println(result)
			self.Output(result)
}

func (self *MirrorHandler) GET() {
	self.Run()
}

func (self *MirrorHandler) POST() {
	self.Run()
}

var app_example *Application

func init() {
	app_example = Init()
	app_example.Route("^/mirror(/(?P<args>[0-9a-zA-Z]*))?$", &MirrorHandler{})
	app_example.Route("^/tpl$", &TplHandler{})
	go app_example.Run(":12345")
}
