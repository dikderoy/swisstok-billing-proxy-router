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
	"github.com/spf13/viper"
	"fmt"
	"time"
)

func mainFatalLogFromRecover(exitCode int) {
	if err := recover(); err != nil {
		fmt.Printf("main.fatal: %s \n", err)
		os.Exit(exitCode)
	}
}

func init() {
	defer mainFatalLogFromRecover(1)
	//defaults
	viper.SetTypeByDefaultValue(true)
	viper.SetDefault("ESBEndpoint", "http://httpbin.org/post?esb=1")
	viper.SetDefault("AVKEndpoint", "http://httpbin.org/post?avk=1")
	viper.SetDefault("bucket-size", 100)

	//config file setup
	viper.SetConfigName("config")
	viper.SetEnvPrefix("har")
	viper.AddConfigPath("/etc/hap/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

func main() {
	defer mainFatalLogFromRecover(2)

	bucketTtl, err := time.ParseDuration(viper.GetString("bucket-ttl"))
	if err != nil {
		panic("failed to parse bucket-ttl, check config")
	}
	models.InitGlobalBucket(bucketTtl)
	a := app.App{}
	a.Init()
	app.Logger = log.New(os.Stdout, "\nhttp:", 0)

	r := mux.NewRouter()
	//get routes configuration
	var endpointConfig map[string]handlers.RequestHandler

	if err := viper.UnmarshalKey("routes", &endpointConfig); err != nil {
		panic(fmt.Sprintf("cant unmarshal routes: %s , check config", err))
	}

	var handler app.Handler
	for route, conf := range endpointConfig {
		switch conf.SenderType {
		case models.SenderESB:
			handler = &handlers.EsbRequestHandler{conf}
			break;
		case models.SenderAVK:
			handler = &handlers.AvkRequestHandler{conf}
			break;
		default:
			panic("unknown route type!")
		}
		fmt.Printf("applying new route: %v\n", handler)
		r.Handle(route, app.NewAppHandler(a.NewAppKernel(handler)))
	}
	r.Handle("/status", app.NewAppHandler(
		a.NewAppKernel(&handlers.ListHandler{})))

	fmt.Println("starting application")
	a.Start()
	fmt.Printf("creating server at [%s]\n", viper.GetString("bind-to"))
	server := http.Server{
		Addr: viper.GetString("bind-to"),
		Handler:gHandlers.CombinedLoggingHandler(os.Stdout,
			gHandlers.RecoveryHandler(gHandlers.RecoveryLogger(app.Logger))(r))}
	fmt.Printf("launching server at [%s]\n", viper.GetString("bind-to"))
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server Error: %s", err)
	}
}
