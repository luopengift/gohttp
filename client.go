package gohttp

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	*http.Request
}

func (self *Request) SetHeader(k, v string) *Request {
	self.Header.Add(k, v)
	return self
}

func (self *Request) Headers(kv map[string]string) *Request {
	for k, v := range kv {
		self.SetHeader(k, v)
	}
	return self
}

type Response struct {
	*http.Response
}

func (self *Response) Code() int {
	return self.StatusCode
}

func (self *Response) Bytes() ([]byte, error) {
	return ioutil.ReadAll(self.Body)
}

func (self *Response) String() string {
	if resp, err := self.Bytes(); err != nil {
		return ""
	} else {
		return string(resp)
	}

}

func NewRequest(method, urlStr string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, url.QueryEscape(urlStr), body)
	return &Request{req}, err
}

type Client struct {
	method     string
	url        string
	path       string
	query      string
	fragment   string
	cookies    map[string]string
	headers    map[string]string
	body       interface{}
	proxy      string
	timeout    int
	retries    int
	verifySsl  bool //true:强制使用https,false:不校验https证书
	keepAlived bool
	transport  http.Transport
}

func NewClient() *Client {
	return &Client{
		method:     "GET",
		cookies:    make(map[string]string),
		headers:    make(map[string]string),
		verifySsl:  false,
		keepAlived: true,
	}
}

//构造request body [interface{} -> io.Reader]
func parseBody(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	bts, err := Bytes(v)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(bts)
	return body, nil
}

func (self *Client) newRequest() (*http.Request, error) {
	u, err := self.newURL()
	if err != nil {
		return nil, err
	}

	body, err := parseBody(self.body)
	if err != nil {
		return nil, err
	}

	req, err := NewRequest(self.method, u.String(), body)
	for k, v := range self.headers {
		req.Header.Set(k, v)
	}
	for k, v := range self.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return req.Request, err
}

func (self *Client) setClient() (*http.Client, error) {
	if self.proxy != "" {
		proxy, err := url.Parse(self.proxy)
		if err != nil {
			return nil, err
		}
		self.transport.Proxy = http.ProxyURL(proxy)
	}
	self.transport.DisableKeepAlives = !self.keepAlived
	self.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !self.verifySsl}
	client := &http.Client{Transport: &self.transport}
	client.Timeout = time.Duration(self.timeout) * time.Second
	return client, nil
}

func (self *Client) Reset() *Client {
	self.method = "GET"
	self.url = ""
	self.path = ""
	self.query = ""
	self.fragment = ""
	self.cookies = make(map[string]string)
	self.headers = make(map[string]string)
	self.body = nil
	self.proxy = ""
	self.timeout = 0
	self.retries = 0
	self.verifySsl = false
	self.keepAlived = true
	self.transport = http.Transport{}
	return self
}

func (self *Client) doReq(method string) (*Response, error) {
	self.method = method
	req, err := self.newRequest()
	if err != nil {
		return nil, err
	}

	client, err := self.setClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	return &Response{resp}, err
}

// 长连接,Default is true
func (self *Client) KeepAlived(used bool) *Client {
	self.keepAlived = used
	return self
}

// 强制使用HTTPS,Default is false
func (self *Client) VerifySSL(used bool) *Client {
	self.verifySsl = used
	return self
}

func (self *Client) URL(urlstr string) *Client {
	self.url = urlstr
	return self
}

func (self *Client) Path(path string) *Client {
	self.path = path
	return self
}

func (self *Client) Proxy(proxy string) *Client {
	self.proxy = proxy
	return self
}

func (self *Client) Timeout(timeout int) *Client {
	self.timeout = timeout
	return self
}

func (self *Client) Cookie(k, v string) *Client {
	self.cookies[k] = v
	return self
}

func (self *Client) Header(k, v string) *Client {
	self.headers[k] = v
	return self
}

func (self *Client) Headers(kv map[string]string) *Client {
	for k, v := range kv {
		self.Header(k, v)
	}
	return self
}

func (self *Client) Body(body interface{}) *Client {
	self.body = body
	return self
}

func (self *Client) Retries(count int) *Client {
	self.retries = count
	return self
}

func (self *Client) newURL() (*url.URL, error) {
	u, err := url.Parse(self.url)
	if err != nil {
		return u, err
	}
	if self.path != "" {
		u.Path = self.path
	}
	if self.query != "" {
		u.RawQuery = self.query
	}
	return u, err
}

func (self *Client) Get() (*Response, error)  { return self.doReq("GET") }
func (self *Client) Post() (*Response, error) { return self.doReq("POST") }
func (self *Client) Head() (*Response, error) { return self.doReq("HEAD") }
func (self *Client) Put() (*Response, error)  { return self.doReq("PUT") }
