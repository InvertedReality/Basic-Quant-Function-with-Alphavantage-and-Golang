package main

import (
	"net/http"
	"time"
)

func main() {
	p("Basic QUANT APP", version(), "started at", config.Address)
	mux := http.NewServeMux()

	//The urls
	mux.HandleFunc("/get/", getdata)

	//Server details
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
