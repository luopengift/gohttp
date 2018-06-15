package gohttp

import (
	"testing"
)

func Test_http(t *testing.T) {
	app := Init()
	app.SetTLS("cert.pem", "key.pem")
	app.RunHttp(":8888")
}
