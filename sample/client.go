package main

import (
    "fmt"
    "github.com/luopengift/gohttp"
)
func main() {
    resp,err := gohttp.NewClient().URL("http://www.baidu.com").Header("Content-Type","application/json;charset=utf-8").Get()
    fmt.Println(resp.String(), err)
}
