package models
import "errors"

var bucket *RequestBucket

type RequestBucket struct {
	bucket map[int]*Request
}

func InitGlobalBucket(queueSize int) *RequestBucket {
	bucket = NewRequestBucket(queueSize)
	return bucket
}

func GlobalBucket() *RequestBucket {
	return bucket
}

func NewRequestBucket(queueSize int) *RequestBucket {
	return &RequestBucket{bucket:make(map[int]*Request, queueSize)}
}

func (self *RequestBucket) add(r *Request) {
	self.bucket[r.id] = r
}

func (self *RequestBucket) find(id int) (*Request, error) {
	if r := self.bucket[id]; r != nil {
		return r, nil
	}
	return nil, errors.New("Request does not exists")
}

func (self *RequestBucket) remove(id int) {
	delete(self.bucket, id)
}