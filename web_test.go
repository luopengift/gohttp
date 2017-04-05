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
	fmt.Println(RouterMap)
	go HttpsRun(&Config{
		Addr:     ":443",
		CertFile: "./server.cert",
		KeyFile:  "./server.key",
	})
	go HttpRun(&Config{Addr: ":80"})
	select {}
}
