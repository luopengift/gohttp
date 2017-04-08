# gohttp
simple http framework for golang

### 模拟tornado

#### 使用说明

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
    self.Output([]byte("world"))
}

func (self *Test) POST() {
    self.Output([]byte("hello"))
}

//启动服务
func main() {
    //绑定路由
    gohttp.RouterRegister("^/(?P<ID>[0-9]*)/(?P<NAME>[a-zA-Z]*)$", &Test{})
    gohttp.HttpRun(&gohttp.Config{
        Addr:     ":9999",
        CertFile: "./server.cert",
        KeyFile:  "./server.key",
    })
}
```






