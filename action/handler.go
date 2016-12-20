package action

import "net/http"

type Handler struct {
	Action Action
}

func (c Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		Body:   []byte{},
		Header: http.Header{},
	}

	_, err := c.Action.Do(r, res)

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
