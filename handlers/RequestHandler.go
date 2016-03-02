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

type EsbRequestHandler struct {
	RequestHandler
}

type AvkRequestHandler struct {
	RequestHandler
}

func (self *EsbRequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) *app.AppError {
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
	response, err := r.RequestAVK(self.TargetEndpoint)
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target", err)
	}
	*self.GetApp().GetChannel() <- r
	app.Logger.Printf("serving:", *r)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.GetId()))
	res.Write(response)
	return nil
}

func (self *AvkRequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) *app.AppError {
	var target string = self.TargetEndpoint
	defer req.Body.Close()
	r := models.NewRequest(*req, self.SenderType, self.RequestType, self.CallbackPath)
	if _, err := r.ParseIdFromXMLRes(); err != nil {
		return app.NewAppError(500, "Failed to parse request")
	}
	fmt.Println("id extracted")
	fmt.Println("bucket query")
	if mr, err := models.GlobalBucket().Find(r.GetId()); err == nil && mr.GetId() == r.GetId() {
		target = mr.CallbackAddress
		fmt.Println("bucket queried")
	}
	//pass avk response to esb
	response, err := r.ProxyRequest(target, *req)
	defer response.Body.Close()
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target " + target, err)
	}
	fmt.Println("body extracting")
	body := models.ReadContentFromRequest(response)
	fmt.Println("body extracted")
	*self.GetApp().GetChannel() <- r
	app.Logger.Printf("serving:", *r)
	res.WriteHeader(response.StatusCode)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.GetId()))
	res.Write(body)
	return nil
}
