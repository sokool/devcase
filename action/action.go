package action

import "net/http"

type Action interface {
	Do(http.ResponseWriter, *http.Request) interface{}
}

type ActionFunc func(http.ResponseWriter, *http.Request) interface{}

func (af ActionFunc) Do(w http.ResponseWriter, r *http.Request) interface{} {
	return af(w, r)
}
