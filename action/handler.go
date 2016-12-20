package action

import "net/http"

type handler struct {
	action Action
}

func (c handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		Body:   []byte{},
		Header: http.Header{},
	}

	_, err := c.action.Do(r, res)

	// handle errors
	switch err.(type) {
	case Error:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	case error:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//rewrite headers
	for k, h := range res.Header {
		for _, v := range h {
			w.Header().Set(k, v)
		}
	}
	w.Write(res.Body)
}
