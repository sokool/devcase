package action

import "net/http"

type caller struct {
	Action Action
}

func (c caller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Action.Do(w, r)
}
