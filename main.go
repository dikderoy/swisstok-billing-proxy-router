package main

import (
	"./models"
	"./handlers"
	"./app"
	"net/http"
	"log"
	"os"
	"github.com/gorilla/mux"
	gHandlers "github.com/gorilla/handlers"
)

func main() {
	models.InitGlobalBucket(100)
	a := app.App{}
	a.Init(app.AppConfig{AVKEndpoint:"http://httpbin.org/post"})

	app.Logger = log.New(os.Stdout, "\nhttp:", 0)

	r := mux.NewRouter()
	r.Handle("/request/esb", a.NewAppKernel(handlers.RequestHandler{}))
	r.Handle("/request/avk", a.NewAppKernel(handlers.RequestHandler{}))
	r.Handle("/status", a.NewAppKernel(handlers.ListHandler))

	server := http.Server{
		Addr: ":8081",
		Handler:gHandlers.RecoveryHandler(
			gHandlers.RecoveryLogger(app.Logger))(r)}

	server.ListenAndServe()
}
