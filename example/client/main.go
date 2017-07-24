package main

import (
	"fmt"
	"github.com/luopengift/gohttp"
)

func main() {
	client := gohttp.NewClient().Url("http://www.baidu.com").Header("Content-Type", "application/json;charset=utf-8")
	for i:=0;i<=4;i++ {

		resp, err := client.Get()
		fmt.Println(fmt.Sprintf("%#v", client))
		fmt.Println(resp.Code(), err)

	}
	fmt.Println("======")
	pool := gohttp.NewClientPool(1,4,10)
	for i:=0;i<=10;i++ {
		go func() {
			conn := pool.Get()
			resp,err := conn.Url("http://www.baidu.com").Header("Content-Type", "application/json;charset=utf-8").Get()
			fmt.Println(fmt.Sprintf("%v,%v,%p",resp.Code(), err,conn))
			pool.Put(conn)
		}()
	}
	select{}
}
