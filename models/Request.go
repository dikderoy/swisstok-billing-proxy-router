package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestAllocationError struct {
	info string
}

func (self RequestAllocationError) Error() string {
	return self.info
}

type Request struct {
	id        int
	req       http.Request
	timestamp time.Time
}

func NewRequest(req http.Request) *Request {
	return &Request{req: req, timestamp: time.Now()}
}

func (self Request) String() string {
	return fmt.Sprintf("Req#%d#t.%s", self.id, self.timestamp)
}

func (self Request) proxyRequest(addr string, req http.Request) (*http.Response, error) {
	return http.Post(addr, "application/json", req.Body)
}

func (self *Request) requestTarget() ([]byte, error) {
	resp, err := self.proxyRequest("http://httpbin.org/post", self.req)
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

func (self Request) parseJsonResponse(body []byte) (id int, err error) {
	var f  interface{}
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

/*
func (self *Request) passResponse(req http.Request) error {

	if addr=="" {
		return errors.New("return path is not defined")
	}
	resp, _ := self.proxyRequest(addr, req)
}
*/
