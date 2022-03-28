package main

import (
	"fmt"
	"go-template/config"
	"go-template/router"
	"net/http"
	"time"
)

func main() {
	s := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%s",
			config.Conf.GetString("host.address"),
			config.Conf.GetString("host.port"),
		),
		Handler:        router.Router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
