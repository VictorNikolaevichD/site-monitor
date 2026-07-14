package checker

import (
	"net/http"
	"time"
)

type Result struct {
	URL                string
	Code               int
	AvailabilityStatus bool
	Error              error
}

type Checker struct {
	client *http.Client
}

func NewChecker() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Checker) Check(url string) Result {
	resp, err := c.client.Get(url)
	if err != nil {
		return Result{URL: url, Error: err}
	}
	defer resp.Body.Close()

	return Result{
		URL:                url,
		Code:               resp.StatusCode,
		AvailabilityStatus: resp.StatusCode == http.StatusOK,
		Error:              nil,
	}
}
