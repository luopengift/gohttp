package gohttp

import (
	"net/http"
)

func (self *HttpHandler) Prepare() {}
func (self *HttpHandler) Finish()  {}
func (self *HttpHandler) GET() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) HEAD() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) POST() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) PUT() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) PATCH() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) DELETE() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) OPTIONS() {
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
