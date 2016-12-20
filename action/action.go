package action

import "net/http"

type Action interface {
	Do(*http.Request, *Response) (interface{}, error)
}

type Response struct {
	Header http.Header
	Body   []byte
}

type ActionFunc func(*http.Request, *Response) (interface{}, error)

func (af ActionFunc) Do(req *http.Request, res *Response) (interface{}, error) {
	return af(req, res)
}
