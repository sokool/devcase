package action

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
)

type Wrapper func(Action) Action

//Decorate allows you to do something with Action request and response
//before or after is executed. For instance you can log, transform response,
//monitor...
func Decorate(a Action, ws ...Wrapper) http.Handler {
	for _, wrapper := range ws {
		a = wrapper(a)
	}

	return Handler{
		Action: a,
	}
}

//JsonResponse takes action result and creates JSON HTTP response
func JsonResponse() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(req *http.Request, res *Response) (interface{}, error) {
			acc := req.Header.Get("accept")
			if acc != "application/xml" {
				req.Header.Set("accept", "application/json")
			}
			out, err := a.Do(req, res)
			if err != nil {
				return out, err
			}

			if req.Header.Get("accept") == "application/json" {
				res.Body, err = json.Marshal(out)
				if err != nil {
					//internal error
					return out, err
				}

				res.Header.Set("content-type", "application/json")
			}
			return out, err
		})
	}
}

//XMLResponse takes action result and creates XML HTTP response
func XMLResponse() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(req *http.Request, res *Response) (interface{}, error) {
			out, err := a.Do(req, res)
			if err == nil && req.Header.Get("accept") == "application/xml" {
				res.Body, err = xml.Marshal(out)
				if err != nil {
					return out, err
				}
				res.Header.Set("content-type", "application/xml")
			}

			return out, err
		})
	}
}

func Logger() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(req *http.Request, res *Response) (interface{}, error) {
			out, err := a.Do(req, res)

			log.Printf("header %s", res.Header)
			log.Printf("body %s", res.Body)
			return out, err
		})
	}
}
