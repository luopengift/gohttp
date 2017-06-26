package gohttp

import (
	//"fmt"
	"testing"
)

func Test_web(t *testing.T) {
	app := Init()
	app.Route("^/mirror(/(?P<args>[0-9a-zA-Z]*))?$", &MirrorHandler{})
	app.Route("^/tpl", &TplHandler{})
	//fmt.Println(app.String())
	//fmt.Println(fmt.Sprintf("%#v", app.Server))
	app.Run()
}

/*
func Test_client(t *testing.T) {
	fmt.Println("test")
	client, err := NewClient().URL("https://www.luopengift.com").Path("/test").Body("Hello").Cookie("appversion", "1.5.0").Post()
	fmt.Println("resp", client.Response)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("response", client)
	}
}
*/
