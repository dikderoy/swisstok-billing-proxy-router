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
	config AppConfig
}

func (self *App) Init() {
	self.chanel = make(chan *models.Request, 1000)
}

func (self *App) Start() {
	go self.HandleIncomingRequests()
}

func (self *App) HandleIncomingRequests() {
	for r := range self.chanel {
		if _, err := models.GlobalBucket().Find(r.GetId()); err == nil || r.GetType() == models.SenderAVK {
			Logger.Printf("touch bucket: %s", *r)
			//models.GlobalBucket().Remove(r.GetId())
			continue
		}
		models.GlobalBucket().Add(r)
		Logger.Printf("add to bucket: %s", *r)
	}
}

func (self *App) Config() AppConfig {
	return self.config
}

func (self *App) GetChannel() *chan *models.Request {
	return &self.chanel
}

func (app *App) NewAppKernel(h Handler) Handler {
	h.SetApp(app)
	return h
}