package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// API Endpoints {Produce, Consume}
/*
Unmarshal the request’s JSON body into a struct.
Run that endpoint’s logic with the request to obtain a result.
Marshal and write that result to the response.
*/

func NewHTTPServer(addr string) *http.Server{
    httpsvr := newHTTPServer()
    r := mux.NewRouter()

    r.handleFunc("/", httpsvr.handleConsume).Methods("GET")
    r.handleFunc("/", httpsvr.handleProduce).Methods("POST")
    return &http.Server{
        Addr: addr,
        Handler: r,
    }
}

type httpServer struct {
    Log *Log
}

func newHTTPServer(){

}
