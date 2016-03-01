package app

import (
	"log"
	"../models"
	"net/http"
	"time"
	"io"
	"strings"
	"strconv"
)

var (
	Logger *log.Logger
)

type App struct {
	chanel chan *models.Request
	bucket *models.RequestBucket
	config AppConfig
}

func (self *App) Init(config AppConfig) {
	self.config = config
	self.chanel = make(chan *models.Request, 1000)
}

func (self *App) Start() {
	go self.HandleIncomingRequests()
}

func (self *App) HandleIncomingRequests() {
	for r := range self.chanel {
		if req, err := models.GlobalBucket().Find(r.GetId()); err == nil {
			Logger.Printf("remove from bucket:", *req)
			models.GlobalBucket().Remove(r.GetId())
		}else {
			models.GlobalBucket().Add(r)
			Logger.Printf("add to bucket:", *r)
		}
	}
}

func (self *App) SimulateResponse(r *models.Request) {
	time.Sleep(10*1000000)
	var buf io.Reader = strings.NewReader(`<?xml version="1.0" encoding="utf-8"?>` +
	`<request><req_type>wait_payment</req_type><param><cobill>17942</cobill><cobillgroup>9928</cobillgroup><corequest_list><corequest>` +
	strconv.Itoa(r.GetId()) + `</corequest></corequest_list><currency>RUB</currency><sum>300</sum></param></request>`)
	Logger.Printf("simulating xml response for", *r)
	http.Post("http://localhost:8081/request/avk", models.RequestTypeXML, buf)
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