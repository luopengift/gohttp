# gohttp
simple http framework for golang

### 模拟tornado

#### 使用说明
1. 启动服务
```
package main

import (
    "github.com/luopengift/gohttp"
)
//Handler方法实现
type Test struct {
    HttpHandler
}

func (self *Test) GET() {
    self.Output([]byte("world"))
}

func (self *Test) POST() {
    self.Output([]byte("hello"))
}

//绑定路由
gohttp.RouterRegister("^/(?P<ID>[0-9]*)/(?P<NAME>[a-zA-Z]*)$", &Test{})
//启动服务
func main() {
    gohttp.Start(&Config{
        Addr:     ":9999",
        certFile: "./server.cert",
        keyFile:  "./server.key",
    })
}
```
