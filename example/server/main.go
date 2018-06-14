package main

import (
	"github.com/luopengift/gohttp/example/assets"
)

func main() {
	assets.App.RunHttp(":18081")
}
