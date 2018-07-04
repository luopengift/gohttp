package gohttp

import (
	//"io/ioutil"
	"net/http"
	"strings"
)

type request struct {
	// native *http.Request
	*http.Request

	// request method
	Method string

	Proto string

	// request url
	URL string

	// request host without port
	Remote string
	//
	Path string
	// request path regx match arguments
	match map[string]string
	body  []byte
}

func NewRequestReader(req *http.Request) *request {
	r := new(request)
	r.Request = req
	r.Method = req.Method
	r.Proto = req.Proto
	r.URL = req.RequestURI
	r.Remote = strings.Split(req.RemoteAddr, ":")[0]
	r.Path = req.URL.Path
	r.match = make(map[string]string)
	return r
}

// parse and handler arguments
func (req *request) parse_arguments(match map[string]string) error {
	// parse form automatically
	if err := req.Request.ParseForm(); err != nil {
		return err
	}
	req.prepare_match_arguments(match)
	return nil
}

// prepare match and assignment to match arguments
func (req *request) prepare_match_arguments(match map[string]string) {
	req.match = match
}

func (req *request) GetCookies() []*http.Cookie {
	return req.Request.Cookies()
}

func (req *request) GetCookie(name string) string {
	cookie, err := req.Request.Cookie(name)
	if err != nil {
		panic(err)
	}
	return cookie.Value
}
