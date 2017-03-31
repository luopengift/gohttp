package gohttp

import (
	"fmt"
	"net/http"
	"time"
)

func HttpRun(config *Config) {
	fmt.Println("HttpServer Start", config.Addr)
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        NewHttpHandler(),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func HttpsRun(config *Config) {
	fmt.Println("HttpsServer Start", config.Addr)
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        NewHttpHandler(),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil {
		fmt.Println(err)
	}
}
