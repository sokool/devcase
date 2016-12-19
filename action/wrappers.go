package action

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Wrapper func(Action) Action

func Wrapper(a Action, ws ...Wrapper) http.Handler {
	for _, wrapper := range ws {
		a = wrapper(a)
	}

	return caller{
		Action: a,
	}
}

func Json() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(w http.ResponseWriter, r *http.Request) interface{} {
			out := a.Do(w, r)
			if r.Header.Get("content-type") == "application/json" {
				json.NewEncoder(w).Encode(out)
			}
			return out
		})
	}
}

func Xml() Wrapper {
	return func(a Action) Action {
		return ActionFunc(func(w http.ResponseWriter, r *http.Request) interface{} {
			out := a.Do(w, r)
			if r.Header.Get("content-type") == "application/xml" {
				if err := xml.NewEncoder(w).Encode(out); err != nil {
					return out
				}
			}
			return out
		})
	}
}
