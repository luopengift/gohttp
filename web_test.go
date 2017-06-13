package gohttp

import (
	"fmt"
	"testing"
)

func Test_web(t *testing.T) {
	fmt.Println(RouterMap)
	/*go HttpsRun(&Config{
		Addr:     ":443",
		CertFile: "./server.cert",
		KeyFile:  "./server.key",
	})
	*/
	HttpRun(&Config{Addr: ":18080"})
	//select {}
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
