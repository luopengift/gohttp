package gohttp

import ()

type Test struct {
	HttpHandler
}

func (self *Test) GET() {
	self.Output([]byte("Match query args-->\n"))
	for k, v := range self.GetMatchArgs() {
		self.Output([]byte(k + ":" + v + "\n"))
	}
	self.Output([]byte("query args-->\n"))
	for k, v := range self.GetQueryArgs() {
		self.Output([]byte(k + ":" + v[0] + "\n"))
	}
	self.Output([]byte(self.GetQueryArg("a")))
}

func (self *Test) POST() {
	self.Output([]byte(self.Request().PostFormValue("id")))
	self.Output([]byte("hello"))
}

func init() {
	RouterRegister("^/(?P<ID>[0-9]*)/(?P<NAME>[a-zA-Z]*)$", &Test{})
}
