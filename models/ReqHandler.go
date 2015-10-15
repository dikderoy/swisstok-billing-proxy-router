package models
import (
	"log"
	"net/http"
	"fmt"
	"bytes"
	"strconv"
	"io/ioutil"
)

type RequestHandler struct {
	requests  map[int]Request
	increment *int
	logger    *log.Logger
}

func ReqHandler(l *log.Logger, queueSize int) RequestHandler {
	counter := 0
	obj := RequestHandler{}
	obj.increment = &counter
	obj.requests = make(map[int]Request, queueSize)
	obj.logger = l

	return obj
}

func (self *RequestHandler) addRequest(r *Request) {
	self.requests[*self.increment] = *r
}

func (self RequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	*self.increment++
	r := Request{id:*self.increment, req:*req}
	length, _ := strconv.ParseInt(req.Header.Get("Content-Length"), 10, 1)
	self.logger.Print("content-length:", length, req.Header.Get("Content-Length"))
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err {
	}
	self.addRequest(&r)
	self.logger.Printf("serving:", *self.increment, bytes.NewBuffer(body).String())

	res.Header().Add("X-Request-Id", fmt.Sprint(*self.increment))
	res.Write(bytes.NewBufferString("Response Data").Bytes())
}