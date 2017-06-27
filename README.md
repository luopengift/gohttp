# gohttp [![GoWalker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/luopengift/gohttp)

gohttp is used for RESTful APIs, Web apps, Http services in Golang.
It is used similar with [Tornado](http://www.tornadoweb.org).


#### GO verion >=1.8.0

#### 使用说明
1. HTTP server
```
package main

import (
    "github.com/luopengift/gohttp"
)

//Handler方法实现
type Test struct {
    gohttp.HttpHandler
}

func (self *Test) GET() {
    self.Output("world")
}

func (self *Test) POST() {
    self.Output("hello")
}

//启动服务
func main() {
    app := gohttp.Init() 
    //绑定路由
    app.Route("^/(?P<ID>[0-9]*)/(?P<NAME>[a-zA-Z]*)$", &Test{})
    app.Run("8080")
}
```
2. HTTP client
```
resp,err := gohttp.NewClient().URL(http://www.google.com).Header("Content-Type","application/json;charset=utf-8").Get()
fmt.Println(resp.String())
```





