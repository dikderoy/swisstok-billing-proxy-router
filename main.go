package main

import (
	"net/http"
	"log"
	"os"
	"./models"
)

func main() {
	logger := log.New(os.Stdout, "http:", 0)
	mux := http.NewServeMux()
	server := http.Server{Addr:":8081",Handler:mux}
	mux.Handle("/request/new",models.ReqHandler(logger,100))
	server.ListenAndServe()
}
