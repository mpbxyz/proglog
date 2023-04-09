package main

import (
	"log"

	"github.com/mpbxyz/proglog/internal/server"
)
func main(){
    svr := server.NewHTTPServer(":8080")
    log.Fatal(svr.ListenAndServe())
}
