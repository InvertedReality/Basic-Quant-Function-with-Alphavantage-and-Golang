package main

import (
	"net/http"
	"time"
	"encoding/json"
    "github.com/gorilla/mux"
)


func main() {
	p("Basic QUANT APP", version(), "started at Adress:", config.Address, "\n", time.Now())
	router := mux.NewRouter()

	//The urls
	router.HandleFunc("/", Home)
	router.HandleFunc("/get/data/", ValidationMiddleware(getdata))
	router.HandleFunc("/get/meta/", ValidationMiddleware(GetMetaData))
	router.HandleFunc("/get/graph/", ValidationMiddleware(getgraph))
	router.HandleFunc("/user/logout/", ValidationMiddleware(authenticate)
	router.HandleFunc("/user/auth/", authenticate)
	router.HandleFunc("/user/signup/",ValidationMiddleware(UserExec))
	router.HandleFunc("/user/list/",ValidationMiddleware(UserExec))
	router.HandleFunc("/user/delete/",ValidationMiddleware(UserExec))
	//Server details
	server := &http.Server{
		Addr:           config.Address,
		Handler:        router,
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
		9: "/token/",
	}

	{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(urls)

	}

}

