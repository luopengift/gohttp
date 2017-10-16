package gohttp

import (
	"testing"
)

func Test_http(t *testing.T) {
	app := Init()
	app.SetAddress(":8888")
	app.Run()
}
