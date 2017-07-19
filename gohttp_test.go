package gohttp

import (
	"testing"
)

func Test_http(t *testing.T) {
	app := Init()
	app.SetAddress(":8088")
	app.Run()
}
