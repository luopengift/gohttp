package gohttp

import (
	"net/http"
	"time"
)

func Start(config *Config) {
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        NewHttpHandler(),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	if config.CertFile != "" && config.KeyFile != "" {
		server.ListenAndServeTLS(config.CertFile, config.KeyFile)
	} else {
		server.ListenAndServe()
	}
}
