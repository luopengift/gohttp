package gohttp

import (
	"net/http"
)

func (self *HttpHandler) Prepare() {}
func (self *HttpHandler) Finish()  {}
func (self *HttpHandler) GET() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) HEAD() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) POST() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) PUT() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) PATCH() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) DELETE() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func (self *HttpHandler) OPTIONS() {
	//If defines GET method,must rewrite this function.
	http.Error(self.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
