package gohttp

import ()


type pprofHandler struct {
    HttpHandler
}

func (self *pprofHandler) GET() {

}





func init() {
    RouterRegister("^/pprof$",&pprofHandler{})
}
