package gohttp

import (
	"net/http"
)

type RequestHandler struct {
    *Conn
    matchArgs map[string]string
    queryArgs map[string][]string
    bodyArgs  []byte
    finished  bool
}

func (self *RequestHandler) RemoteAddr() string {
	return self.Request.RemoteAddr
}

func (self *RequestHandler) Prepare() {}
func (self *RequestHandler) Finish()  {}
func (self *RequestHandler) GET() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *RequestHandler) HEAD() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *RequestHandler) POST() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *RequestHandler) PUT() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *RequestHandler) PATCH() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *RequestHandler) DELETE() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *RequestHandler) OPTIONS() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
