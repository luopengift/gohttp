package gohttp

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	*http.Request
}

func (req *Request) SetHeader(k, v string) *Request {
	req.Header.Add(k, v)
	return req
}

func (req *Request) SetHeaders(kv map[string]string) *Request {
	for k, v := range kv {
		req.SetHeader(k, v)
	}
	return req
}

type Response struct {
	*http.Response
}

func (resp *Response) Code() int {
	return resp.StatusCode
}

func (resp *Response) Bytes() ([]byte, error) {
	return ioutil.ReadAll(resp.Body)
}

func (resp *Response) String() string {
	if response, err := resp.Bytes(); err != nil {
		return ""
	} else {
		return string(response)
	}

}

func NewRequest(method, urlStr string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
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
	bts, err := ToBytes(v)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(bts)
	return body, nil
}

func (c *Client) newRequest() (*http.Request, error) {
	u, err := c.newURL()
	if err != nil {
		return nil, err
	}

	body, err := parseBody(c.body)
	if err != nil {
		return nil, err
	}

	req, err := NewRequest(c.method, u.String(), body)
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	for k, v := range c.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return req.Request, err
}

func (c *Client) setClient() (*http.Client, error) {
	if c.proxy != "" {
		proxy, err := url.Parse(c.proxy)
		if err != nil {
			return nil, err
		}
		c.transport.Proxy = http.ProxyURL(proxy)
	}
	c.transport.DisableKeepAlives = !c.keepAlived
	c.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !c.verifySsl}
	client := &http.Client{Transport: &c.transport}
	client.Timeout = time.Duration(c.timeout) * time.Second
	return client, nil
}

func (c *Client) Reset() *Client {
	c.method = "GET"
	c.url = ""
	c.path = ""
	c.query = ""
	c.fragment = ""
	c.cookies = make(map[string]string)
	c.headers = make(map[string]string)
	c.body = nil
	c.proxy = ""
	c.timeout = 0
	c.retries = 0
	c.verifySsl = false
	c.keepAlived = true
	c.transport = http.Transport{}
	return c
}

func (c *Client) doReq(method string) (*Response, error) {
	c.method = method
	req, err := c.newRequest()
	if err != nil {
		return nil, err
	}

	client, err := c.setClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	return &Response{resp}, err
}

// 长连接,Default is true
func (c *Client) KeepAlived(used bool) *Client {
	c.keepAlived = used
	return c
}

// 强制使用HTTPS,Default is false
func (c *Client) VerifySSL(used bool) *Client {
	c.verifySsl = used
	return c
}

func (c *Client) URL(urlstr string) *Client {
	c.url = urlstr
	return c
}

func (c *Client) Path(path string) *Client {
	c.path += path
	return c
}

func (c *Client) Query(kv map[string]string) *Client {
	query := []string{}
	for k, v := range kv {
		s := k + "=" + url.QueryEscape(v)
		query = append(query, s)
	}
	c.query = strings.Join(query, "&")
	return c
}

func (c *Client) Proxy(proxy string) *Client {
	c.proxy = proxy
	return c
}

func (c *Client) Timeout(timeout int) *Client {
	c.timeout = timeout
	return c
}

func (c *Client) Cookie(k, v string) *Client {
	c.cookies[k] = v
	return c
}

func (c *Client) Header(k, v string) *Client {
	c.headers[k] = v
	return c
}

func (c *Client) Headers(kv map[string]string) *Client {
	for k, v := range kv {
		c.Header(k, v)
	}
	return c
}

func (c *Client) Body(body interface{}) *Client {
	c.body = body
	return c
}

func (c *Client) Retries(count int) *Client {
	c.retries = count
	return c
}

func (c *Client) newURL() (*url.URL, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return u, err
	}
	if c.path != "" {
		u.Path = c.path
	}
	if c.query != "" {
		u.RawQuery = c.query
	}
	return u, err
}

func (c *Client) URLString() string {
	url, _ := c.newURL()
	return url.String()
}

func (c *Client) Get() (*Response, error)  { return c.doReq("GET") }
func (c *Client) Post() (*Response, error) { return c.doReq("POST") }
func (c *Client) Head() (*Response, error) { return c.doReq("HEAD") }
func (c *Client) Put() (*Response, error)  { return c.doReq("PUT") }
