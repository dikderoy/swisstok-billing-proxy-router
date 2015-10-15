package models

import (
	"fmt"
	"log"
	"net/http"
)

var Requests map[int]*Request

type RequestHandler struct {
	logger *log.Logger
}

func NewRequestHandler(l *log.Logger, queueSize int) RequestHandler {
	obj := RequestHandler{}
	obj.logger = l
	Requests = make(map[int]*Request, queueSize)
	return obj
}

func (self *RequestHandler) addRequest(r *Request) {
	Requests[r.id] = r
}

func (self RequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	defer func() { if r := recover(); r != nil {
		http.Error(res, "Recovered from EFatal", 500)
	} }()

	r := NewRequest(*req)
	response, err := r.requestBilling()
	if err != nil {
		self.logger.Printf("%#v", r.id);
		http.Error(res, err.Error(), 500)
		return
	}
	self.addRequest(&r)
	self.logger.Printf("serving:", r)
	res.Header().Add("X-Request-Id", fmt.Sprint(r.id))
	res.Write(response)
}
