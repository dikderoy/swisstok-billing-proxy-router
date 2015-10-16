package main

import (
	"./models"
	"net/http"
	"./app"
	"log"
	"os"
)

func main() {
	models.InitGlobalBucket(100)
	app.Logger = log.New(os.Stdout, "\nhttp:", 0)


	mux := http.NewServeMux()
	server := http.Server{Addr: ":8081", Handler: mux}
	mux.Handle("/request/new", app.AppHandler(models.RequestHandler))
	mux.Handle("/status", app.AppHandler(models.ListHandler))
	server.ListenAndServe()
}
