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
	//	RequestTypeXML = "application/xml"
	//	RequestTypeJSON = "application/json"

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
		req:             req,
		timestamp:       time.Now(),
		sender:          sender,
		CallbackAddress: cb,
		cType:           cType}
}

func (self *Request) GetId() int {
	return self.id
}

func (self *Request) GetType() string {
	return self.sender
}

func (self Request) String() string {
	return fmt.Sprintf("Req#%d#t.%s", self.id, self.timestamp)
}

func (self *Request) ProxyRequest(addr string, req http.Request) (response *http.Response, err error) {
	fmt.Printf("\nproxy request{%d} to %s\n", self.id, addr)
	response, err = http.Post(addr, self.cType, req.Body)
	fmt.Printf("\nproxy request complited with status: %3d %s\n", response.StatusCode, response.Status)
	return
}

func (self *Request) ParseIdFromJSON(resp *http.Response) (body []byte, err error) {
	var response AccessibleHttpResponse
	response = AccessibleHttpResponse(*resp)
	body = ReadContentFromRequest(&response)
	self.id, err = self.ParseJsonResponse(body)
	if err != nil {
		return body, &RequestAllocationError{code:1, info:fmt.Sprintf("Ex:json.traverseId: %s", err)}
	}
	return body, nil
}

func (self *Request) ParseIdFromXML() ([]byte, error) {
	var err error
	fmt.Println("parse xml response")
	var request AccessibleHttpRequest
	request = AccessibleHttpRequest(self.req)
	body := ReadContentFromRequest(&request)
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
	//slice root
	ujs = ujs.Get("result")
	fmt.Println(*ujs)
	var idNode *simplejson.Json
	idNode, ok := ujs.Get("data").GetIndex(0).CheckGet("COREQUEST");
	if !ok {
		if idNode, ok = ujs.Get("values").CheckGet("v_corequest"); !ok {
			if idNode, ok = ujs.CheckGet("result"); !ok {
				return 0, &RequestAllocationError{code:1, info:"failed parse request"}
			}
		}
	}
	if sid, err := idNode.String(); err == nil {
		if id_f64, err := strconv.ParseFloat(sid, 64); err == nil {
			id = int(id_f64)
		}
	}
	return
}

func (self Request) ParseXmlResponse(body []byte) (id int, err error) {
	var f []interface{}
	bReader := bytes.NewReader(body)
	f, err = x2j.ReaderValuesForTag(bReader, "corequest")
	if err != nil || len(f) == 0 {
		fmt.Print("len is 0")
		return 0, &RequestAllocationError{code:1, info:fmt.Sprintf("Ex:xml.traverseId: %s", err)}
	}
	fmt.Println("no errors catched", err, len(f), f)
	fid, err := strconv.ParseFloat(f[:1][0].(string), 64)
	id = int(fid)
	return
}

type HttpMessageAbstract interface {
	readBody() (body []byte, err error)
}

type AccessibleHttpResponse http.Response

func (resp *AccessibleHttpResponse) readBody() (body []byte, err error) {
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return
}

type AccessibleHttpRequest http.Request

func (req *AccessibleHttpRequest) readBody() (body []byte, err error) {
	body, err = ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return
}

func ReadContentFromRequest(httpAbstract HttpMessageAbstract) (bodyBytes []byte) {
	bodyBytes, err := httpAbstract.readBody()
	if err != nil {
		panic("failed to get body")
	}
	fmt.Println("\n--req_body begin--")
	bytes.NewBuffer(bodyBytes).WriteTo(os.Stdout)
	fmt.Println("\n--req_body end--")
	return
}
