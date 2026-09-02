package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChecker_Check(t *testing.T) {
	testTable := []struct {
		name          string
		setup         func(t *testing.T) string
		timeout       time.Duration
		wantCode      int
		wantAvailable bool
		wantErr       bool
	}{
		{
			name: "returns OK on 200",
			setup: func(t *testing.T) string {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				t.Cleanup(server.Close)
				return server.URL
			},
			timeout:       2 * time.Second,
			wantCode:      http.StatusOK,
			wantAvailable: true,
		},
		{
			name: "returns unavailable on 404",
			setup: func(t *testing.T) string {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
				t.Cleanup(server.Close)
				return server.URL
			},
			timeout:       2 * time.Second,
			wantCode:      http.StatusNotFound,
			wantAvailable: false,
		},
		{
			name: "returns unavailable on 503",
			setup: func(t *testing.T) string {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}))
				t.Cleanup(server.Close)
				return server.URL
			},
			timeout:       2 * time.Second,
			wantCode:      http.StatusServiceUnavailable,
			wantAvailable: false,
		},
		{
			name: "returns error on connection timeout",
			setup: func(t *testing.T) string {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(500 * time.Millisecond)
					w.WriteHeader(http.StatusOK)
				}))
				t.Cleanup(server.Close)
				return server.URL
			},
			timeout:       50 * time.Millisecond,
			wantAvailable: false,
			wantErr:       true,
		},
		{
			name: "returns error when server is unavailable",
			setup: func(t *testing.T) string {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				url := server.URL
				server.Close()
				return url
			},
			timeout:       2 * time.Second,
			wantAvailable: false,
			wantErr:       true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			url := testCase.setup(t)

			checker := NewChecker(testCase.timeout)
			result := checker.Check(url)

			if result.URL != url {
				t.Fatalf("URL = %q, want %q", result.URL, url)
			}
			if (result.Error != nil) != testCase.wantErr {
				t.Fatalf("Error = %v, wantErr = %v", result.Error, testCase.wantErr)
			}
			if result.Code != testCase.wantCode {
				t.Fatalf("Code = %d, want %d", result.Code, testCase.wantCode)
			}
			if result.AvailabilityStatus != testCase.wantAvailable {
				t.Fatalf("AvailabilityStatus = %v, want %v", result.AvailabilityStatus, testCase.wantAvailable)
			}
			if result.Duration < 0 {
				t.Fatalf("Duration = %v, want >= 0", result.Duration)
			}
		})
	}
}
