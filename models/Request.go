package models

import (
	"fmt"
	"github.com/beego/x2j"
	"io/ioutil"
	"net/http"
	"time"
	"bytes"
	"strconv"
	"github.com/bitly/go-simplejson"
	"os"
)

const (
	RequestTypeXML = "application/xml"
	RequestTypeJSON = "application/json"

	SenderAVK = "AVK"
	SenderESB = "ESB"
)

type RequestAllocationError struct {
	code int
	info string
}

func (self RequestAllocationError) Code() int {
	return self.code
}

func (self RequestAllocationError) Error() string {
	return self.info
}

type Request struct {
	id              int
	req             http.Request
	timestamp       time.Time
	sender          string
	cType           string
	CallbackAddress string
}

func NewRequest(req http.Request, sender, cType string, cb string) *Request {
	return &Request{
		req: req,
		timestamp: time.Now(),
		sender: sender,
		CallbackAddress: cb,
		cType: cType}
}

func (self *Request) GetId() int {
	return self.id
}

func (self Request) String() string {
	return fmt.Sprintf("Req#%d#t.%s", self.id, self.timestamp)
}

func (self *Request) ProxyRequest(addr string, req http.Request) (*http.Response, error) {
	fmt.Printf("proxy request{%d} to %s", self.id, addr)
	return http.Post(addr, self.cType, req.Body)
}

func (self *Request) RequestAVK(addr string) ([]byte, error) {
	resp, err := self.ProxyRequest(addr, self.req)
	defer resp.Body.Close()
	body := ReadContentFromRequest(resp)
	self.id, err = self.ParseJsonResponse(body)
	if err != nil {
		return body, RequestAllocationError{code:1, info:"failed to extract id"}
	}
	return body, nil
}

func (self *Request) ParseIdFromXML() ([]byte, error) {
	var err error
	fmt.Println("parse xml response")
	body := ReadContentFromRequest(&self.req)
	self.id, err = self.ParseXmlResponse(body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (self *Request) ParseJsonResponse(body []byte) (id int, err error) {
	ujs, err := simplejson.NewJson(body)
	if err != nil {
		return
	}
	//slice test prefix
	//ujs = ujs.Get("json")
	if id, err = ujs.GetPath("result", "data").GetIndex(0).Get("COREQUEST").Int(); err == nil {
		return
	}
	return ujs.GetPath("result", "result").Int()
}

func (self Request) ParseXmlResponse(body []byte) (id int, err error) {
	var f []interface{}
	bReader := bytes.NewReader(body)
	f, err = x2j.ReaderValuesFromTagPath(bReader, "request.param.corequest_list.corequest")
	if err != nil {
		return 0, RequestAllocationError{code:1, info:"Ex:xml.traverseId:" + err.Error()}
	}
	fmt.Println(f)
	fid, err := strconv.ParseFloat(f[:1][0].(string), 64)
	id = int(fid)
	return
}

func ReadContentFromRequest(httpAbstract interface{}) (bodyBytes []byte) {
	var err error = nil
	if y, ok := httpAbstract.(*http.Response); ok == true {
		bodyBytes, err = ioutil.ReadAll(y.Body)
		y.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	if x, ok := httpAbstract.(*http.Request); ok == true {
		bodyBytes, err = ioutil.ReadAll(x.Body)
		x.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	if err != nil || len(bodyBytes) == 0 {
		panic("failed to get body")
	}
	fmt.Println("--req_body begin--")
	bytes.NewBuffer(bodyBytes).WriteTo(os.Stdout)
	fmt.Println("--req_body end--")
	return
}