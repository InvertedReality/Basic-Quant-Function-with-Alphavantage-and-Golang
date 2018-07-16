package main

import (
	"net/http"
	"time"
)

func main() {
	p("Basic QUANT APP", version(), "started at Adress:", config.Address, "\n", time.Now())
	mux := http.NewServeMux()

	//The urls
	mux.HandleFunc("/get/data/", getdata)
	mux.HandleFunc("/get/meta/", GetMetaData)
	mux.HandleFunc("/get/graph/", getgraph)
	mux.HandleFunc("/logout/", logout)
	mux.HandleFunc("/signup/", signupAccount)
	mux.HandleFunc("/authenticate/", authenticate)
	mux.HandleFunc("/", Home)

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
