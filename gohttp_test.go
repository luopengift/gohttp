package gohttp

import (
	"testing"
)

func Test_http(t *testing.T) {
	app := Init()
	app.SetTLS("cert.pem", "key.pem")
	app.Run(":8888")
}
