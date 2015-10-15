package models
import "net/http"

type Request struct {
	id int
	req http.Request
	res http.Response
}

func (self *Request) forgeResponse() {

}

func (self *Request) forgeRequest() {

}