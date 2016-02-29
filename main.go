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
	a.Init(app.AppConfig{
		AVKEndpoint:"http://httpbin.org/post",
		ESBEndpoint:"http://httpbin.org/post",
	})

	app.Logger = log.New(os.Stdout, "\nhttp:", 0)

	r := mux.NewRouter()
	r.Handle("/request/esb", app.NewAppHandler(
		a.NewAppKernel(&handlers.RequestHandler{
			TargetEndpoint:a.Config().AVKEndpoint,
			RequestType:models.RequestTypeJSON,
			SenderType:models.SenderESB,
		})))
	r.Handle("/request/avk", app.NewAppHandler(
		a.NewAppKernel(&handlers.RequestHandler{
			TargetEndpoint:a.Config().ESBEndpoint,
			RequestType:models.RequestTypeXML,
			SenderType:models.SenderAVK,
			CallbackPath:"http://localhost:8081",
		})))
	r.Handle("/status", app.NewAppHandler(
		a.NewAppKernel(&handlers.ListHandler{})))

	a.Start()

	server := http.Server{
		Addr: ":8081",
		Handler:gHandlers.RecoveryHandler(
			gHandlers.RecoveryLogger(app.Logger))(r)}

	server.ListenAndServe()
}
