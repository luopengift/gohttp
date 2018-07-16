package gohttp

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/luopengift/golibs/pool"
	"github.com/luopengift/log"
	"github.com/luopengift/types"
)

// ClientPool cilent pool
type ClientPool struct {
	*pool.Pool
}

// NewClientPool new cilent pool
func NewClientPool(maxIdle, maxOpen, timeout int) *ClientPool {
	client := func() (interface{}, error) {
		return NewClient().Reset(), nil
	}
	p := pool.NewPool(maxIdle, maxOpen, timeout, client)
	return &ClientPool{Pool: p}
}

// Get signel client
func (p *ClientPool) Get() (*Client, error) {
	one, err := p.Pool.Get()
	if err != nil {
		log.Error("Get Client error:%v", err)
		return nil, err
	}
	return one.(*Client), nil
}

// Put one client into pool
func (p *ClientPool) Put(c *Client) error {
	return p.Pool.Put(c)
}

// Client client
type Client struct {
	*http.Client
	*http.Transport
	*url.URL
	body    io.Reader
	headers map[string]string
	cookies map[string]string
}

// NewClient new client
func NewClient() *Client {
	c := new(Client)
	c.Client = new(http.Client)
	c.Transport = new(http.Transport)
	c.Transport.TLSClientConfig = new(tls.Config)
	c.Client.Transport = c.Transport
	return c.Reset()
}

// Reset reset
func (c *Client) Reset() *Client {
	c.URL = new(url.URL)
	c.body = nil
	c.headers = make(map[string]string)
	c.cookies = make(map[string]string)
	return c
}

// KeepAlived ,Default is true
func (c *Client) KeepAlived(used bool) *Client {
	c.Transport.DisableKeepAlives = used
	return c
}

// VerifySSL 强制使用HTTPS,Default is false
func (c *Client) VerifySSL(used bool) *Client {
	c.Transport.TLSClientConfig.InsecureSkipVerify = !used
	return c
}

// Timeout timeout
func (c *Client) Timeout(timeout int) *Client {
	c.Client.Timeout = time.Duration(timeout) * time.Second
	return c
}

// Proxy proxy
func (c *Client) Proxy(proxy string) *Client {
	_proxy, err := url.Parse(proxy)
	if err != nil {
		log.Error("proxy set fail:%v", err)
		return c
	}
	c.Transport.Proxy = http.ProxyURL(_proxy)
	return c

}

// URLString url string
func (c *Client) URLString(urlstr string) *Client {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Error("Url set fail:%v", err)
		return c

	}
	//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	c.URL.Scheme = u.Scheme
	c.URL.Opaque = u.Opaque         // encoded opaque data
	c.URL.User = u.User             // username and password information
	c.URL.Host = u.Host             // host or host:port
	c.URL.Path = u.Path             // path (relative paths may omit leading slash)
	c.URL.RawPath = u.RawPath       // encoded path hint (see EscapedPath method)
	c.URL.ForceQuery = u.ForceQuery // append a query ('?') even if RawQuery is empty
	c.URL.RawQuery = u.RawQuery     // encoded query values, without '?
	c.URL.Fragment = u.Fragment     // fragment for references, without '#'
	return c
}

// Path path
func (c *Client) Path(path string) *Client {
	c.URL.Path = path
	return c
}

// Query query params
func (c *Client) Query(query map[string]string) *Client {
	q := []string{}
	for k, v := range query {
		q = append(q, k+"="+url.QueryEscape(v))
	}
	c.URL.RawQuery = strings.Join(q, "&")
	return c
}

// Body body params
func (c *Client) Body(v interface{}) *Client {
	if v == nil {
		return c
	}
	bts, err := types.ToBytes(v)
	if err != nil {
		log.Error("body set fail:%v", err)
		return c
	}
	c.body = bytes.NewBuffer(bts)
	return c
}

// Header header
func (c *Client) Header(k, v string) *Client {
	c.headers[k] = v
	return c
}

// Headers headers
func (c *Client) Headers(kv map[string]string) *Client {
	for k, v := range kv {
		c.headers[k] = v
	}
	return c
}

// Cookie cookie
func (c *Client) Cookie(k, v string) *Client {
	c.cookies[k] = v
	return c
}

// BaseAuth base auth
func (c *Client) BaseAuth(user, pass string) *Client {
	s := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	c.Header("Authorization", "Basic "+s)
	return c
}

func (c *Client) doReq(method string) (*Response, error) {

	req, err := http.NewRequest(method, c.URL.String(), c.body)
	if err != nil {
		log.Error("new request fail:%v", err)
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
		log.Error("client do fail:%v", err)
		return nil, err
	}

	response, err := NewResponse(resp)
	if err != nil {
		log.Error("response read fail:%v", err)
		return nil, err
	}
	return response, nil

}

// Response response
type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0
	Byte       []byte
}

// NewResponse new response
func NewResponse(resp *http.Response) (*Response, error) {
	response := new(Response)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response.Status = resp.Status
	response.StatusCode = resp.StatusCode
	response.Proto = resp.Proto
	response.ProtoMajor = resp.ProtoMajor
	response.ProtoMinor = resp.ProtoMinor
	response.Byte = body
	return response, nil
}

// Code status code
func (resp *Response) Code() int {
	return resp.StatusCode
}

func (resp *Response) String() string {
	return string(resp.Byte)
}

// Bytes bytes
func (resp *Response) Bytes() []byte {
	return resp.Byte
}

//Get request
func (c *Client) Get() (*Response, error) { return c.doReq("GET") }

// Post request
func (c *Client) Post() (*Response, error) { return c.doReq("POST") }

// Put request
func (c *Client) Put() (*Response, error) { return c.doReq("PUT") }

//Delete request
func (c *Client) Delete() (*Response, error) { return c.doReq("DELETE") }

// Head request
func (c *Client) Head() (*Response, error) { return c.doReq("HEAD") }
