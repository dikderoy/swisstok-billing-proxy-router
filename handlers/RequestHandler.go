package handlers

import (
	"fmt"
	"net/http"
	"../app"
	"../models"
)

type RequestHandler struct {
	app.AppKernel
	app.AppHandler
}

func (self *RequestHandler) Handle(res http.ResponseWriter, req *http.Request) *app.AppError {
	defer req.Body.Close()
	addr := req.Header.Get("X-Callback-Path")
	if addr == "" {
		return app.NewAppError(400, "Callback path not given")
	}
	r := models.NewRequest(*req, models.SenderESB, models.RequestTypeJSON)
	response, err := r.RequestTarget(self.GetApp().Config().AVKEndpoint)
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target")
	}
	models.GlobalBucket().Add(r)
	app.Logger.Printf("serving:", *r)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.GetId()))
	res.Write(response)
	return nil
}
