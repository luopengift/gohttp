package gohttp

import (
	"net/http"
)

type Connection interface {
    Request() *http.Request
    Response() http.ResponseWriter
}

type Conn struct {
	request  *http.Request
	response http.ResponseWriter
}

func (self *Conn) init(response http.ResponseWriter, request *http.Request) *Conn {
	self.response = response
	self.request = request
	return self
}

func (self *Conn) Response() http.ResponseWriter {
	return self.response
}

func (self *Conn) Request() *http.Request {
	return self.request
}

