package gohttp

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func Test_web(t *testing.T) {
	fmt.Println(os.Getwd())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("start...")
	fmt.Println(Router)
	go Start(&Config{
		Addr:     ":443",
		CertFile: "./server.cert",
		KeyFile:  "./server.key",
	})
	go Start(&Config{Addr: ":80"})
	select {}
}
