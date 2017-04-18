package gohttp

import (
	"fmt"
/*
	"os"
	"runtime"
*/
	"testing"
)
/*
func Test_web(t *testing.T) {
	fmt.Println(os.Getwd())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(RouterMap)
	go HttpsRun(&Config{
		Addr:     ":443",
		CertFile: "./server.cert",
		KeyFile:  "./server.key",
	})
	go HttpRun(&Config{Addr: ":8080"})
	select {}
}*/

func Test_client(t *testing.T) {
    fmt.Println("test")
    client,err := NewClient().URL("https://www.luopengift.com").Path("/test").Body("Hello").Cookie("appversion","1.5.0").Post()
    fmt.Println("resp",client.Response)
    if err != nil {
        fmt.Println(err)
    }else{
        fmt.Println("response",client)
    }
}
