package models

import (
	"errors"
	"github.com/ikoroteev/ttlcache"
	"time"
	"strconv"
)

var bucket *RequestBucket

type RequestBucket struct {
	bucket *ttlcache.Cache
	ttl    time.Duration
}

func InitGlobalBucket(ttl time.Duration) *RequestBucket {
	bucket = NewRequestBucket(ttl)
	return bucket
}

func GlobalBucket() *RequestBucket {
	return bucket
}

func (self *RequestBucket) Count() int {
	return self.bucket.Count()
}

func NewRequestBucket(ttl time.Duration) *RequestBucket {
	return &RequestBucket{bucket:ttlcache.NewCache(), ttl:ttl}
}

func (self *RequestBucket) Add(r *Request) {
	self.bucket.Set(convert(r.GetId()), r, self.ttl)
}

func (self *RequestBucket) Find(id int) (*Request, error) {
	if r, found := self.bucket.Get(convert(id), true); found == true {
		return r.(*Request), nil
	}
	return nil, errors.New("Request does not exists")
}

func (self *RequestBucket) Remove(id int) {
	self.bucket.Delete(convert(id))
}

func convert(in int) string {
	return strconv.Itoa(in)
}