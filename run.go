package gohttp

import (
	"fmt"
	"net/http"
    "time"
)

func Start(config *Config) {
	fmt.Println(Router)
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        NewHttpHandler(),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
    if config.certFile != "" && config.keyFile != "" {
        server.ListenAndServeTLS(config.certFile,config.keyFile)
    }else{
        server.ListenAndServe()
    }
}
