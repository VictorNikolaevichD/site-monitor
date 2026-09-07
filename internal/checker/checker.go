package checker

import (
	"io"
	"net/http"
	"time"
)

type Result struct {
	URL                string
	Code               int
	AvailabilityStatus bool
	Duration           time.Duration
	Error              error
}

type Checker struct {
	client *http.Client
}

func NewChecker(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Checker) Check(url string) Result {
	startedAt := time.Now()

	resp, err := c.client.Get(url)
	duration := time.Since(startedAt)

	if err != nil {
		return Result{URL: url, Duration: duration, Error: err}
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	return Result{
		URL:                url,
		Duration:           duration,
		Code:               resp.StatusCode,
		AvailabilityStatus: resp.StatusCode < http.StatusBadRequest,
		Error:              nil,
	}
}
