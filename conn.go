package gohttp

import (
	"net/http"
)

type Connection interface {
	request() *http.Request
	response() *http.Response
	responseWriter() http.ResponseWriter
}

type Conn struct {
	*http.Request
	*http.Response
	http.ResponseWriter
}

func (self *Conn) init(responsewriter http.ResponseWriter, request *http.Request) *Conn {
	self.ResponseWriter = responsewriter
	self.Request = request
	return self
}

func (self *Conn) responseWriter() http.ResponseWriter {
	return self.ResponseWriter
}

func (self *Conn) request() *http.Request {
	return self.Request
}

func (self *Conn) response() *http.Response {
	return self.Response
}
