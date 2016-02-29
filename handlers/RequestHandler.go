package handlers

import (
	"fmt"
	"net/http"
	"../app"
	"../models"
)

type RequestHandler struct {
	app.Kernel
	SenderType     string
	TargetEndpoint string
	RequestType    string
	CallbackPath   string
}

func (self *RequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) *app.AppError {
	var cbPath string
	defer req.Body.Close()
	addr := req.Header.Get("X-Callback-Path")
	if addr != "" {
		cbPath = addr
	} else if self.CallbackPath != "" {
		cbPath = self.CallbackPath
	} else {
		return app.NewAppError(400, "Callback path not given")
	}
	r := models.NewRequest(*req, self.SenderType, self.RequestType, cbPath)
	response, err := r.RequestTarget(self.TargetEndpoint)
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target", err)
	}
	*self.GetApp().GetChannel() <- r
	app.Logger.Printf("serving:", *r)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.GetId()))
	res.Write(response)
	return nil
}
