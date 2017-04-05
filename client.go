package gohttp

import (
    "net/http"
    "io"
    "fmt"
)


type Request struct {
    *http.Request
}



func NewRequest(method,urlStr string, body io.Reader) *Request {
    req ,err := http.NewRequest(method, urlStr, body)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    return &Request{req}
}


func (self *Request) Do() (*http.Response,error) {
    client := &http.Client{}
    return client.Do(self.Request)
}


