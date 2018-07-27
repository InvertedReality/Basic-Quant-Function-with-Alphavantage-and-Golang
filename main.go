package main

import (
	"net/http"
	"time"
	"encoding/json"
	// "fmt"
	// "reflect"
)

func main() {
	p("Basic QUANT APP", version(), "started at Adress:", config.Address, "\n", time.Now())
	mux := http.NewServeMux()

	//The urls
	mux.HandleFunc("/get/data/", getdata)
	mux.HandleFunc("/get/meta/", GetMetaData)
	mux.HandleFunc("/get/graph/", getgraph)
	mux.HandleFunc("/user/logout/", authenticate)
	mux.HandleFunc("/user/signup/",UserExec)
	mux.HandleFunc("/user/auth/", authenticate)
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/user/list/",UserExec)
	mux.HandleFunc("/user/delete/",UserExec)

	//Server details
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	return
}

//home
//lists urls
func Home(writer http.ResponseWriter, request *http.Request) {
	urls := map[int]string{
		1: "/get/data",
		2: "/get/meta",
		3: "/get/graph",
		4: "/user/logout",
		5: "/user/signup",
		6: "/user/auth/",
		7: "/user/list/",
		8: "/user/delete/",
	}

	{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(urls)

	}

}

