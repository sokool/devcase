package action

import "net/http"

type Action interface {
	Do(http.ResponseWriter, *http.Request) (interface{}, error)
}

type ActionFunc func(http.ResponseWriter, *http.Request) (interface{}, error)

func (af ActionFunc) Do(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return af(w, r)
}
