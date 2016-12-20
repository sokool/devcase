package fixer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Base  string
	Date  string
	Rates map[string]float64
}

var (
	host string = "http://api.fixer.io/"
)

func Stream(base string, d time.Duration) (<-chan Response, error) {
	c := make(chan Response)
	t := time.NewTicker(d)

	first, err := Latest(base)
	if err != nil {
		return c, err
	}

	go func(cr chan<- Response) {
		cr <- first
		for range t.C {
			r, err := Latest(base)
			if err != nil {
				log.Printf("could't connect to %s due %v", host, err)
				continue
			}
			cr <- r
		}
	}(c)

	return c, nil
}

func Latest(base string) (Response, error) {
	return fetch("latest", map[string]string{"base": base})
}

func fetch(endpoint string, query map[string]string) (Response, error) {
	var s Response
	var q string = ""
	for n, r := range query {
		q += fmt.Sprintf("%s=%s&", n, r)
	}

	r, err := http.Get(fmt.Sprintf("%s/%s?%s", host, endpoint, q))
	if err != nil {
		return s, err
	}

	if r.StatusCode >= http.StatusBadRequest {
		o, _ := ioutil.ReadAll(r.Body)
		return s, fmt.Errorf("status code %s [%v]", r.Status, string(o))
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return s, err
	}

	return newResponse(b)
}

func newResponse(b []byte) (Response, error) {
	r := Response{
		Rates: make(map[string]float64),
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return r, err
	}

	return r, nil
}
