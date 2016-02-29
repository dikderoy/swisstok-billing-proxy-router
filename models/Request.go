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

func (self *Request) proxyRequest(addr string, req http.Request) (*http.Response, error) {
	fmt.Printf("proxy request{%d} to %s", self.id, addr)
	return http.Post(addr, self.cType, req.Body)
}

func (self *Request) RequestTarget(addr string) ([]byte, error) {
	switch self.sender {
	case SenderESB:
		return self.requestAVK(addr)
	case SenderAVK:
		return self.requestESB(addr)
	}
	return []byte{}, RequestAllocationError{"no proxy target"}
}

func (self *Request) requestAVK(addr string) ([]byte, error) {
	resp, err := self.proxyRequest(addr, self.req)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	self.id, err = self.parseJsonResponse(body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (self *Request) requestESB(addr string) ([]byte, error) {
	fmt.Println("read xml response")
	body, err := ioutil.ReadAll(self.req.Body)

	bodyS := bytes.NewBuffer(body).String()
	fmt.Println("rbody:", bodyS)

	if err != nil {
		return []byte{}, err
	}
	fmt.Println("parse xml response")
	self.id, err = self.parseXmlResponse(body)
	if err != nil {
		return body, err
	}
	//here fake req to esb.
	return body, nil
}

func (self *Request) parseJsonResponse(body []byte) (id int, err error) {
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

func (self Request) parseXmlResponse(body []byte) (id int, err error) {
	var f []interface{}
	fmt.Printf("parsing xml response")

	//xmlDoc, err := x2j.ByteDocToMap(body)
	/*
	if err != nil {
		return 0, RequestAllocationError{"Ex:xml.parse:" + err.Error()}
	}
	fmt.Printf("mapping xml response")
	f, err = x2j.MapValue(xmlDoc, "request.id", nil)*/
	breader := bytes.NewReader(body)
	f, err = x2j.ReaderValuesFromTagPath(breader, "request.param.corequest_list.corequest")
	if err != nil {
		return 0, RequestAllocationError{"Ex:xml.traverseId:" + err.Error()}
	}
	fmt.Println(f)

	fid, err := strconv.ParseFloat(f[:1][0].(string), 64)
	id = int(fid)
	return

	//return 0, RequestAllocationError{"xml.id wasn't catched - cant match request"}
}