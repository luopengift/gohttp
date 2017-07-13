package main

import (
    "fmt"
    "github.com/luopengift/gohttp"
)
func main() {
    client := gohttp.NewClient().URL("http://www.baidu.com").Header("Content-Type","application/json;charset=utf-8")
	for _,_ = range []int{1,2,3,4} {
		
		resp,err := client.Get()
		fmt.Println(fmt.Sprintf("%#v",client))	
    		fmt.Println(resp.Code(), err)
	}

}
