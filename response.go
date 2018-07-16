package gohttp

import (
	"bufio"
	"net"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	// ResponseWriter have three method:
	// Header() Header <1>
	// Write([]byte) (int, error) <2>
	// WriteHeader(int) <3>
	http.ResponseWriter

	http.Flusher
	http.Hijacker
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int
	// Finished returns whether or not the ResponseWriter has been finished.
	Finished() bool
	// Size returns the size of the response body.
	Size() int
}

func newResponseWriter(responsewriter http.ResponseWriter) *response {
	resp := new(response)
	resp.ResponseWriter = responsewriter
	return resp
}

// type response implements http.ResponseWriter, http.Hijacker, http.CloseNotify, http.Flush interface
type response struct {
	http.ResponseWriter
	status int
	size   int
}

// rewrite http.ResponseWriter interface method
func (resp *response) Header() http.Header {
	return resp.ResponseWriter.Header()
}

// rewrite http.ResponseWriter interface method
func (resp *response) Write(b []byte) (int, error) {
	size, err := resp.ResponseWriter.Write(b)
	resp.size += size
	return size, err
}

// rewrite http.ResponseWriter interface method
func (resp *response) WriteHeader(code int) {
	resp.status = code
	resp.ResponseWriter.WriteHeader(code)
}

func (resp *response) Status() int {
	return resp.status
}

func (resp *response) Size() int {
	return resp.size
}

func (resp *response) Finished() bool {
	return resp.status != 0
}

// Implements the http.Hijacker interface
func (resp *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if resp.size < 0 {
		resp.size = 0
	}
	return resp.ResponseWriter.(http.Hijacker).Hijack()
}

// Implements the http.CloseNotify interface
func (resp *response) CloseNotify() <-chan bool {
	return resp.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Implements the http.Flush interface
func (resp *response) Flush() {
	resp.ResponseWriter.(http.Flusher).Flush()
}
