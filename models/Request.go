package models

import (
	"encoding/json"
	"fmt"
	"github.com/beego/x2j"
	"io/ioutil"
	"net/http"
	"time"
	"bytes"
	"strconv"
	"go/types"
)

const (
	RequestTypeXML = "application/xml"
	RequestTypeJSON = "application/json"

	SenderAVK = "AVK"
	SenderESB = "ESB"
)

type RequestAllocationError struct {
	info string
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
	callbackAddress string
}

func NewRequest(req http.Request, sender, cType string, cb string) *Request {
	return &Request{
		req: req,
		timestamp: time.Now(),
		sender: sender,
		callbackAddress: cb,
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

func (self *Request) Process(addr string) ([]byte, error) {
	switch self.sender {
	case SenderESB:
		return self.RequestAVK(addr)
	case SenderAVK:
		return self.RequestESB(addr)
	}
	return []byte{}, RequestAllocationError{"no proxy target"}
}

func (self *Request) RequestAVK(addr string) ([]byte, error) {
	resp, err := self.ProxyRequest(addr, self.req)
	body := ReadContentFromRequest(resp)
	self.id, err = self.ParseJsonResponse(body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (self *Request) RequestESB(addr string) ([]byte, error) {
	var err error
	fmt.Println("parse xml response")
	body := ReadContentFromRequest(self.req)
	self.id, err = self.ParseXmlResponse(body)
	if err != nil {
		return body, err
	}
	//here fake req to esb.
	return body, nil
}

func (self *Request) ParseJsonResponse(body []byte) (id int, err error) {
	var f interface{}
	if err = json.Unmarshal(body, &f); err != nil {
		return
	}
	for k, v := range f.(map[string]interface{}) {
		if k == "json" {
			for k2, v2 := range v.(map[string]interface{}) {
				if k2 == "id" {
					switch v2.(type) {
					case float64:
						id = int(v2.(float64))
						return
					default:
						break
					}
				}
			}
			break
		}
	}
	return 0, RequestAllocationError{"id wasn't catched - cant allocate request"}
}

func (self Request) ParseXmlResponse(body []byte) (id int, err error) {
	var f []interface{}
	bReader := bytes.NewReader(body)
	f, err = x2j.ReaderValuesFromTagPath(bReader, "request.param.corequest_list.corequest")
	if err != nil {
		return 0, RequestAllocationError{"Ex:xml.traverseId:" + err.Error()}
	}
	fmt.Println(f)
	fid, err := strconv.ParseFloat(f[:1][0].(string), 64)
	id = int(fid)
	return
}

func ReadContentFromRequest(req interface{}) []byte {
	var bodyBytes []byte

	//todo: something wrong here

	switch req.(type) {
	case http.Request:
		bodyBytes, _ = ioutil.ReadAll(req.(http.Request).Body)
		var x http.Request = req.(http.Request)
		x.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		break
	case http.Response:
		bodyBytes, _ = ioutil.ReadAll(req.(http.Response).Body)
		var y http.Response = req.(http.Response)
		y.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		break
	default:
		panic("given instance is nor http.Request, nor http.Response")
	}
	return bodyBytes
}