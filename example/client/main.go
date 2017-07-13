package main

import (
	"fmt"
	"github.com/luopengift/gohttp"
)

func main() {
	client := gohttp.NewClient().Url("http://www.baidu.com").Header("Content-Type", "application/json;charset=utf-8")
	for {

		resp, err := client.Get()
		fmt.Println(fmt.Sprintf("%#v", client))
		fmt.Println(resp.Code(), err)
		
	}

}
