package handlers

import (
	"fmt"
	"net/http"
	"../app"
	"../models"
	"bytes"
	"io/ioutil"
)

type RequestHandler struct {
	app.Kernel
	SenderType     string `mapstructure:"type"`
	TargetEndpoint string `mapstructure:"target-endpoint"`
	RequestType    string `mapstructure:"content-type"`
	CallbackPath   string `mapstructure:"default-callback-path"`
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
	response, err := r.ProxyRequest(self.TargetEndpoint, *req)
	defer response.Body.Close()
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target", err)
	}
	body, err := r.ParseIdFromJSON(response)
	if err != nil {
		app.Logger.Printf("parsing error: %s", err)
	}
	if r.GetId() > 0 {
		*self.GetApp().GetChannel() <- r
	}
	app.Logger.Printf("serving: %s \n", *r)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.GetId()))
	res.Write(body)
	return nil
}

func (self *AvkRequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) *app.AppError {
	app.Logger.Println("\n-----rh begin-----")
	defer app.Logger.Println("\n-----rh end-----")
	//default target
	var target string = self.TargetEndpoint
	defer req.Body.Close()
	r := models.NewRequest(*req, self.SenderType, self.RequestType, self.CallbackPath)
	if body, err := r.ParseIdFromXML(); err == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		//if valid id - search db
		app.Logger.Println("id extracted")
		app.Logger.Println("bucket query")
		if mr, err := models.GlobalBucket().Find(r.GetId()); err == nil && mr.GetId() == r.GetId() {
			target = mr.CallbackAddress
			app.Logger.Println("bucket queried")
		}
	} else {
		app.Logger.Printf("parsing error: %s", err)
	}
	//pass avk response to esb
	response, err := r.ProxyRequest(target, *req)
	defer response.Body.Close()
	if err != nil {
		return app.NewAppError(500, "Failed to proxy request to target " + target, err)
	}
	body := models.ReadContentFromRequest(response)
	if r.GetId() > 0 {
		//store only identifiable requests
		*self.GetApp().GetChannel() <- r
	}
	app.Logger.Printf("serving: %s \n", *r)
	res.WriteHeader(response.StatusCode)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.GetId()))
	res.Write(body)
	return nil
}
