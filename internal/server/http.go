package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// API Endpoints {Produce, Consume}
/*
Unmarshal the request’s JSON body into a struct.
Run that endpoint’s logic with the request to obtain a result.
Marshal and write that result to the response.
*/

func NewHTTPServer(addr string) *http.Server {
	httpsvr := newHTTPServer()
	r := mux.NewRouter()

	r.HandleFunc("/", httpsvr.handleConsume).Methods("GET")
	r.HandleFunc("/", httpsvr.handleProduce).Methods("POST")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest

	// Decode, append to log, produce output, encode into writer

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(&res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest

	// Decode, append to log, produce output, encode into writer

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rec, err := s.Log.Read(req.Offset)

	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: rec}
	err = json.NewEncoder(w).Encode(&res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}
