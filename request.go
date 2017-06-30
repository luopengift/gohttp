package gohttp

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type request struct {
	// native *http.Request
	*http.Request

	// request method
	Method string

	// request url
	URL string

	// request host without port
	Remote string
	//
	Path string
	// request path regx match arguments
	match map[string]string
	//request query arguments
	query map[string][]string
	// request body arguments
	body []byte
	// TODO:request form arguments
	form map[string][]string
}

func NewRequestReader(req *http.Request) *request {
	r := new(request)
	r.Request = req
	r.Method = req.Method
	r.URL = req.RequestURI
	r.Remote = strings.Split(req.RemoteAddr, ":")[0]
	r.Path = req.URL.Path
	r.match = make(map[string]string)
	r.query = make(map[string][]string)
	r.body = []byte{}
	r.form = make(map[string][]string)
	return r
}

// parse and handler arguments
func (req *request) parse_arguments(match map[string]string) {
	// parse form automatically
	req.Request.ParseForm()

	req.prepare_match_arguments(match)
	req.prepare_query_arguments()
	req.prepare_body_arguments()
	req.prepare_form_arguments()
	//logger.Debug("header:%#v", req.Request.Header)
	//logger.Debug("match:%#v,query:%#v,body:%#v", req.match, req.query, req.body)
	//logger.Debug("PostForm:%#v,MultipartForm:%#v", req.Request.PostForm, req.Request.MultipartForm)
}

// prepare match and assignment to match arguments
func (req *request) prepare_match_arguments(match map[string]string) {
	req.match = match
}

// prepare query and assignment to query arguments
func (req *request) prepare_query_arguments() {
	req.query = req.Request.Form
}

// prepare body and assignment to body arguments
func (req *request) prepare_body_arguments() {
	var err error
	req.body, err = ioutil.ReadAll(req.Request.Body)
	if err != nil {
		panic(err)
	}
}

// prepare form and assignment to form arguments
// Content-Type:application/x-www-form-urlencoded
func (req *request) prepare_form_arguments() {
	req.form = req.Request.PostForm
}
