package main

import (
	"./models"
	"log"
	"net/http"
	"os"
	"./app"
)

func main() {
	app.Logger = log.New(os.Stdout, "\nhttp:", 0)
	mux := http.NewServeMux()
	server := http.Server{Addr: ":8081", Handler: mux}
	mux.Handle("/request/new", models.NewRequestHandler(app.Logger, 100))
	mux.Handle("/status", app.AppHandler(models.ListHandler))
	server.ListenAndServe()
}
