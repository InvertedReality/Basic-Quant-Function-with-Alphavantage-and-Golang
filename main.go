package main

import (
	"net/http"
	"time"
	"encoding/json"
    "github.com/gorilla/mux"
)

var path_details []string

func main() {
	p("Basic QUANT APP", version(), "started at Adress:", config.Address, "\n", time.Now())

    router := mux.NewRouter()

	//The urls
	router.HandleFunc("/", Home).Name("Home")
	router.HandleFunc("/get/data/", ValidationMiddleware(getdata)).Name("data")
	router.HandleFunc("/get/meta/", ValidationMiddleware(GetMetaData)).Name("meta_data")
	router.HandleFunc("/get/graph/", ValidationMiddleware(getgraph)).Name("get_graph")
	router.HandleFunc("/user/logout/", ValidationMiddleware(authenticate)).Name("logout")
	router.HandleFunc("/user/auth/", authenticate).Name("authenticate")
	router.HandleFunc("/user/signup/",UserExec).Name("signup")
	router.HandleFunc("/user/list/",ValidationMiddleware(UserExec)).Name("user_list")
	router.HandleFunc("/user/delete/",ValidationMiddleware(UserExec)).Name("user_delete")
    router.Walk(WalkFunc)

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
	{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(path_details)
	}

}

func WalkFunc(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
    url,_:=route.URLPath()
    path_details= append(path_details,url.Path)
    return nil
}
