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

var client = &http.Client{
	Timeout: time.Second * 10,
}

func CheckSite(url string) Result {
	resp, err := client.Get(url)
	if err != nil {
		return Result{
			Error: err,
		}
	}

	return Result{
		URL:                url,
		Code:               resp.StatusCode,
		AvailabilityStatus: resp.StatusCode == 200,
		Error:              nil,
	}
}
