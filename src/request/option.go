package request

import (
	"github.com/go-resty/resty/v2"
	"time"
)

type Request struct {
	client     *resty.Client
	activityID int
}

type Option func(r *Request)

func (opt Option) Apply(r *Request) {
	opt(r)
}

func New(opts ...Option) *Request {
	r := &Request{
		client: resty.New().SetTimeout(15 * time.Second).
			SetHeaders(map[string]string{
				"user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
				"referer":    "https://www.bilibili.com",
			}),
	}

	for _, o := range opts {
		o.Apply(r)
	}

	return r
}

//func WithDB(d *db.DB) Option {
//	return func(r *Request) {
//		r.db = d
//	}
//}

func WithActivityID(id int) Option {
	return func(r *Request) {
		r.activityID = id
	}
}
