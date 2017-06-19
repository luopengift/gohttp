package main

import (
    "fmt"
    "github.com/luopengift/gohttp"
)

func main() {
    fmt.Println(gohttp.RouterMap)
    gohttp.HttpRun(gohttp.NewConfig(":9995", 100, 100, 100))

}
