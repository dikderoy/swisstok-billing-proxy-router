package app

import (
	"log"
	"../models"
)

var (
	Logger *log.Logger
)

type App struct {
	chanel chan *models.Request
	bucket *models.RequestBucket
	config *AppConfig
}

func (self *App) Init(config *AppConfig) {
	self.config = config
	self.chanel = make(chan *models.Request, 1000)
	self.bucket = models.GlobalBucket()
}

func (self *App) Start() {
	go self.HandleIncomingRequests()
}

func (self *App) HandleIncomingRequests() {
	for r := range App.chanel {
		App.bucket.Add(r)
	}
}

func (self *App) Config() *AppConfig {
	return self.config
}

func (app *App) NewAppKernel(h AppHandler) *AppKernel {
	h.SetApp(app)
	return &AppKernel(h)
}