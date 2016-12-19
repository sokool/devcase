package action

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Wrapper func(Action) Action

func Decorate(a Action, ws ...Wrapper) http.Handler {
	for _, wrapper := range ws {
		a = wrapper(a)
	}

	return caller{
		Action: a,
	}
}

func JsonResponse(d bool) Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
			if d && r.Header.Get("content-type") == "" {
				r.Header.Set("content-type", "application/json")
			}

			out, err := a.Do(w, r)
			if r.Header.Get("content-type") == "application/json" {
				json.NewEncoder(w).Encode(out)
			}
			return out, err
		})
	}
}

func XMLResponse() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
			out, err := a.Do(w, r)
			if r.Header.Get("content-type") == "application/xml" {
				if err := xml.NewEncoder(w).Encode(out); err != nil {
					return out, err
				}
			}
			return out, err
		})
	}
}

func Logger() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
			fmt.Println("before")
			out, err := a.Do(w, r)
			fmt.Println("after")
			return out, err
		})
	}
}
