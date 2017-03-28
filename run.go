package gohttp

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func Start(config *Config) {
	fmt.Println(Router)
	path, _ := exec.LookPath(os.Args[0])
	fmt.Println(path)
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        NewHttpHandler(),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	if config.certFile != "" && config.keyFile != "" {
		server.ListenAndServeTLS(config.certFile, config.keyFile)
	} else {
		server.ListenAndServe()
	}
}
