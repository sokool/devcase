package convert

import (
	"net/http"

	"fmt"
	"strconv"
	"time"

	"github.com/sokool/devcase/action"
	"github.com/sokool/devcase/fixer"
)

var Action action.Action = newConverter()

const Base = "USD"

type result struct {
	Amount    string
	Currency  string
	Converted map[string]string
}

type converterAction struct {
	data  fixer.Response
	error error
}

func (c *converterAction) Do(req *http.Request, res *action.Response) (interface{}, error) {
	currency := req.URL.Query().Get("currency")
	amount := req.URL.Query().Get("amount")

	at, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		return nil, action.Error(fmt.Errorf("error wrong amount format %s, please use ie 10.4321", amount))
	}

	o, err := c.transform(currency, at)

	return o, action.Error(err)
}

// Transform
func (c *converterAction) transform(currency string, amount float64) (result, error) {
	o := result{Converted: map[string]string{}}

	if c.error != nil {
		return o, c.error
	}

	if amount <= 0 {
		return o, fmt.Errorf("error, amount %.2f should not be less/equal than zero", amount)
	}

	o.Amount = fmt.Sprintf("%.2f", amount)
	o.Currency = currency

	if currency == c.data.Base {
		for n, t := range c.data.Rates {
			o.Converted[n] = fmt.Sprintf("%.2f", t)
		}
		return o, nil
	}

	av, ok := c.data.Rates[currency]
	if !ok {
		return o, fmt.Errorf("wrong currency type %s", currency)
	}

	o.Converted[c.data.Base] = fmt.Sprintf("%.2f", amount/av)
	for n, t := range c.data.Rates {
		o.Converted[n] = fmt.Sprintf("%.2f", amount/av*t)
	}

	return o, nil
}

func newConverter() action.Action {
	c := &converterAction{}

	//stream rates from fixer, new stream arrive in every 24hours.
	//this reduce amount of calls to thirt part API
	go func() {
		fc, err := fixer.Stream(Base, time.Hour*24)
		if err != nil {
			c.error = err
			return
		}

		for r := range fc {
			c.data = r
			c.error = nil
		}
	}()

	return c
}
