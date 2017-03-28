package gohttp

import (
    "fmt"
    "runtime"
    "testing"
)

func Test_web(t *testing.T) {
     runtime.GOMAXPROCS(runtime.NumCPU())
    fmt.Println("start...")
    Start(&Config{
        Addr:":9999",
        certFile:   "./server.cert",
        keyFile:    "./server.key",
    })
}
