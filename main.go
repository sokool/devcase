package main

import (
	"net/http"

	"github.com/sokool/devcase/action"
	"github.com/sokool/devcase/convert"
)

func main() {

	http.Handle("/convert", action.Decorate(convert.Action,
		action.XMLResponse(),
		action.JsonResponse(),
		action.Logger(),
	))

	if err := http.ListenAndServe(":9999", nil); err != nil {
		panic(err)
	}

}
