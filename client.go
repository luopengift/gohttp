package gohttp

import (
	"bytes"
	"crypto/tls"
	"github.com/luopengift/golibs/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	*http.Client
	*http.Transport
	*url.URL
	body    io.Reader
	headers map[string]string
	cookies map[string]string
}

func NewClient() *Client {
	c := new(Client)
	c.Client = new(http.Client)
	c.Transport = new(http.Transport)
	c.Transport.TLSClientConfig = new(tls.Config)
	c.Client.Transport = c.Transport
	c.URL = new(url.URL)
	c.body = nil
	c.headers = make(map[string]string)
	c.cookies = make(map[string]string)
	return c
}

// 长连接,Default is true
func (c *Client) KeepAlived(used bool) *Client {
	c.Transport.DisableKeepAlives = used
	return c
}

// 强制使用HTTPS,Default is false
func (c *Client) VerifySSL(used bool) *Client {
	c.Transport.TLSClientConfig.InsecureSkipVerify = !used
	return c
}

func (c *Client) Timeout(timeout int) *Client {
	c.Client.Timeout = time.Duration(timeout) * time.Second
	return c
}

func (c *Client) Proxy(proxy string) *Client {
	_proxy, err := url.Parse(proxy)
	if err != nil {
		logger.Error("proxy set fail:%v", err)
		return c
	}
	c.Transport.Proxy = http.ProxyURL(_proxy)
	return c

}

func (c *Client) Url(urlstr string) *Client {
	u, err := url.Parse(urlstr)
	if err != nil {
		logger.Error("Url set fail:%v", err)
		return c

	}
	//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	c.URL.Scheme = u.Scheme
	c.URL.Opaque = u.Opaque         // encoded opaque data
	c.URL.User = u.User             // username and password information
	c.URL.Host = u.Host             // host or host:port
	c.URL.Path += u.Path            // path (relative paths may omit leading slash)
	c.URL.RawPath = u.RawPath       // encoded path hint (see EscapedPath method)
	c.URL.ForceQuery = u.ForceQuery // append a query ('?') even if RawQuery is empty
	c.URL.RawQuery = u.RawQuery     // encoded query values, without '?
	c.URL.Fragment = u.Fragment     // fragment for references, without '#'
	return c
}

func (c *Client) Path(path string) *Client {
	c.URL.Path += path
	return c
}

func (c *Client) Query(query map[string]string) *Client {
	q := []string{}
	for k, v := range query {
		q = append(q, k+"="+url.QueryEscape(v))
	}
	c.URL.RawQuery = strings.Join(q, "&")
	return c
}

func (c *Client) Body(v interface{}) *Client {
	if v == nil {
		return c
	}
	bts, err := ToBytes(v)
	if err != nil {
		logger.Error("body set fail:%v", err)
		return c
	}
	c.body = bytes.NewBuffer(bts)
	return c
}

func (c *Client) Header(k, v string) *Client {
	c.headers[k] = v
	return c
}

func (c *Client) Headers(kv map[string]string) *Client {
	for k, v := range kv {
		c.headers[k] = v
	}
	return c
}

func (c *Client) Cookie(k, v string) *Client {
	c.cookies[k] = v
	return c
}

func (c *Client) doReq(method string) (*Response, error) {

	req, err := http.NewRequest(method, c.URL.String(), c.body)
	if err != nil {
		logger.Error("new request fail:%v", err)
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	for k, v := range c.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		logger.Error("client do fail:%v", err)
		return nil, err
	}

	defer resp.Body.Close()

	response, err := NewResponse(resp)
	if err != nil {
		logger.Error("response read fail:%v", err)
		return nil, err
	}
	return response, nil

}

type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0
	Byte       []byte
}

func NewResponse(resp *http.Response) (*Response, error) {
	response := new(Response)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response.Status = resp.Status
	response.StatusCode = resp.StatusCode
	response.Proto = resp.Proto
	response.ProtoMajor = resp.ProtoMajor
	response.ProtoMinor = resp.ProtoMinor
	response.Byte = body
	return response, nil
}

func (resp *Response) Code() int {
	return resp.StatusCode
}

func (resp *Response) String() string {
	return string(resp.Byte)
}

func (resp *Response) Bytes() []byte {
	return resp.Byte
}

func (c *Client) Get() (*Response, error)    { return c.doReq("GET") }
func (c *Client) Post() (*Response, error)   { return c.doReq("POST") }
func (c *Client) Put() (*Response, error)    { return c.doReq("PUT") }
func (c *Client) Delete() (*Response, error) { return c.doReq("DELETE") }
func (c *Client) Head() (*Response, error)   { return c.doReq("HEAD") }
