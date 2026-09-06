package site

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site/mocks"
)

var exampleUUID, _ = uuid.Parse("550e8400-e29b-41d4-a716-446655440000")

func TestHandler_GetAll(t *testing.T) {
	testTable := []struct {
		name                      string
		expectedSites             []domain.Site
		executeErr                error
		expectedCallCount         int
		expectedContentTypeHeader string
		expectedStatusCode        int
		expectedResponseBody      string
	}{
		{
			name: "OK",
			expectedSites: []domain.Site{
				{
					ID:   exampleUUID,
					URL:  "https://example.com",
					Name: "Example",
				},
			},
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusOK,
			expectedResponseBody:      "[{\"id\":\"550e8400-e29b-41d4-a716-446655440000\",\"url\":\"https://example.com\",\"name\":\"Example\"}]\n",
		},
		{
			name:                      "OK empty list",
			expectedSites:             []domain.Site{},
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusOK,
			expectedResponseBody:      "[]\n",
		},
		{
			name:                      "Error storage",
			executeErr:                domain.ErrStorage,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusInternalServerError,
			expectedResponseBody:      "{\"error\":\"внутренняя ошибка сервера\"}\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := mocks.NewController()
			c.Return("Execute", testCase.expectedSites, testCase.executeErr)
			getAllUseCase := mocks.NewMockGetAllUseCase(c)

			handler := &Handler{
				getAllUseCase: &getAllUseCase,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/sites", nil)

			handler.GetAll(w, req)

			if got := getAllUseCase.CallCount("Execute"); got != testCase.expectedCallCount {
				t.Fatalf("Expected %d, received %d", testCase.expectedCallCount, got)
			}
			if header := w.Header().Get("Content-Type"); testCase.expectedContentTypeHeader != header {
				t.Fatalf("Expected %s, received %s", testCase.expectedContentTypeHeader, header)
			}
			if testCase.expectedStatusCode != w.Code {
				t.Fatalf("Expected %d, received %d", testCase.expectedStatusCode, w.Code)
			}
			if body := w.Body.String(); testCase.expectedResponseBody != body {
				t.Fatalf("Expected %s, received %s", testCase.expectedResponseBody, body)
			}
		})
	}
}

func TestHandler_Add(t *testing.T) {
	validRequest := "{\"url\":\"https://example.com\",\"name\":\"Example\"}"
	createdSite := domain.Site{
		ID:   exampleUUID,
		URL:  "https://example.com",
		Name: "Example",
	}

	testTable := []struct {
		name                      string
		request                   string
		expectedSite              domain.Site
		executeErr                error
		expectedCallCount         int
		expectedContentTypeHeader string
		expectedStatusCode        int
		expectedResponseBody      string
	}{
		{
			name:                      "OK",
			request:                   validRequest,
			expectedSite:              createdSite,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusCreated,
			expectedResponseBody:      "{\"id\":\"550e8400-e29b-41d4-a716-446655440000\",\"url\":\"https://example.com\",\"name\":\"Example\"}\n",
		},
		{
			name:                      "Error invalid JSON",
			request:                   "{\"bad_field_name\":\"Example\"}",
			expectedCallCount:         0,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"некорректный JSON\"}\n",
		},
		{
			name:                      "Error malformed JSON",
			request:                   "{",
			expectedCallCount:         0,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"некорректный JSON\"}\n",
		},
		{
			name:                      "Error URL required",
			request:                   validRequest,
			executeErr:                domain.ErrURLRequired,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"URL сайта обязателен\"}\n",
		},
		{
			name:                      "Error invalid URL",
			request:                   validRequest,
			executeErr:                domain.ErrInvalidURL,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"некорректный URL сайта\"}\n",
		},
		{
			name:                      "Error name required",
			request:                   validRequest,
			executeErr:                domain.ErrNameRequired,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"имя сайта обязательно\"}\n",
		},
		{
			name:                      "Error site already exists",
			request:                   validRequest,
			executeErr:                domain.ErrSiteAlreadyExists,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusConflict,
			expectedResponseBody:      "{\"error\":\"сайт с таким URL уже существует\"}\n",
		},
		{
			name:                      "Error storage",
			request:                   validRequest,
			executeErr:                domain.ErrStorage,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusInternalServerError,
			expectedResponseBody:      "{\"error\":\"внутренняя ошибка сервера\"}\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := mocks.NewController()
			c.Return("Execute", testCase.expectedSite, testCase.executeErr)
			addUseCase := mocks.NewMockAddUseCase(c)

			handler := &Handler{
				addUseCase: &addUseCase,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sites", bytes.NewBufferString(testCase.request))

			handler.Add(w, req)

			if got := addUseCase.CallCount("Execute"); got != testCase.expectedCallCount {
				t.Fatalf("Expected %d, received %d", testCase.expectedCallCount, got)
			}
			if header := w.Header().Get("Content-Type"); testCase.expectedContentTypeHeader != header {
				t.Fatalf("Expected %s, received %s", testCase.expectedContentTypeHeader, header)
			}
			if testCase.expectedStatusCode != w.Code {
				t.Fatalf("Expected %d, received %d", testCase.expectedStatusCode, w.Code)
			}
			if body := w.Body.String(); testCase.expectedResponseBody != body {
				t.Fatalf("Expected %s, received %s", testCase.expectedResponseBody, body)
			}
		})
	}
}

func TestHandler_DeleteByID(t *testing.T) {
	testTable := []struct {
		name                      string
		id                        string
		executeErr                error
		expectedCallCount         int
		expectedContentTypeHeader string
		expectedStatusCode        int
		expectedResponseBody      string
	}{
		{
			name:                      "OK",
			id:                        exampleUUID.String(),
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusNoContent,
			expectedResponseBody:      "",
		},
		{
			name:                      "Error invalid ID",
			id:                        "not-a-uuid",
			expectedCallCount:         0,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"некорректный ID сайта\"}\n",
		},
		{
			name:                      "Error site not found",
			id:                        exampleUUID.String(),
			executeErr:                domain.ErrSiteNotFound,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusNotFound,
			expectedResponseBody:      "{\"error\":\"сайт с указанным ID не найден\"}\n",
		},
		{
			name:                      "Error storage",
			id:                        exampleUUID.String(),
			executeErr:                domain.ErrStorage,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusInternalServerError,
			expectedResponseBody:      "{\"error\":\"внутренняя ошибка сервера\"}\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := mocks.NewController()
			c.Return("Execute", testCase.executeErr)
			deleteUseCase := mocks.NewMockDeleteByIDUseCase(c)

			handler := &Handler{
				deleteUseCase: &deleteUseCase,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/sites/"+testCase.id, nil)
			req.SetPathValue("id", testCase.id)

			handler.DeleteByID(w, req)

			if got := deleteUseCase.CallCount("Execute"); got != testCase.expectedCallCount {
				t.Fatalf("Expected %d, received %d", testCase.expectedCallCount, got)
			}
			if header := w.Header().Get("Content-Type"); testCase.expectedContentTypeHeader != header {
				t.Fatalf("Expected %s, received %s", testCase.expectedContentTypeHeader, header)
			}
			if testCase.expectedStatusCode != w.Code {
				t.Fatalf("Expected %d, received %d", testCase.expectedStatusCode, w.Code)
			}
			if body := w.Body.String(); testCase.expectedResponseBody != body {
				t.Fatalf("Expected %s, received %s", testCase.expectedResponseBody, body)
			}
		})
	}
}

func TestHandler_GetStatusByID(t *testing.T) {
	checkedAt := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)

	siteWithoutCheck := domain.Site{
		ID:   exampleUUID,
		URL:  "https://example.com",
		Name: "Example",
	}
	siteWithCheck := domain.Site{
		ID:   exampleUUID,
		URL:  "https://example.com",
		Name: "Example",
		LastCheck: &domain.CheckStatus{
			Availability: domain.Available,
			Code:         200,
			CheckedAt:    checkedAt,
			Duration:     245 * time.Millisecond,
		},
	}
	siteUnavailable := domain.Site{
		ID:   exampleUUID,
		URL:  "https://example.com",
		Name: "Example",
		LastCheck: &domain.CheckStatus{
			Availability: domain.Unavailable,
			CheckedAt:    checkedAt,
			Duration:     5 * time.Second,
			Error:        "context deadline exceeded",
		},
	}

	testTable := []struct {
		name                      string
		id                        string
		expectedSite              domain.Site
		executeErr                error
		expectedCallCount         int
		expectedContentTypeHeader string
		expectedStatusCode        int
		expectedResponseBody      string
	}{
		{
			name:                      "OK",
			id:                        exampleUUID.String(),
			expectedSite:              siteWithCheck,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusOK,
			expectedResponseBody:      "{\"url\":\"https://example.com\",\"availability\":\"available\",\"code\":200,\"checked_at\":\"2026-07-15T10:00:00Z\",\"duration_ms\":245}\n",
		},
		{
			name:                      "OK unknown status",
			id:                        exampleUUID.String(),
			expectedSite:              siteWithoutCheck,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusOK,
			expectedResponseBody:      "{\"url\":\"https://example.com\",\"availability\":\"unknown\"}\n",
		},
		{
			name:                      "OK unavailable",
			id:                        exampleUUID.String(),
			expectedSite:              siteUnavailable,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusOK,
			expectedResponseBody:      "{\"url\":\"https://example.com\",\"availability\":\"unavailable\",\"checked_at\":\"2026-07-15T10:00:00Z\",\"duration_ms\":5000,\"error\":\"context deadline exceeded\"}\n",
		},
		{
			name:                      "Error invalid ID",
			id:                        "not-a-uuid",
			expectedCallCount:         0,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusBadRequest,
			expectedResponseBody:      "{\"error\":\"некорректный ID сайта\"}\n",
		},
		{
			name:                      "Error site not found",
			id:                        exampleUUID.String(),
			executeErr:                domain.ErrSiteNotFound,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusNotFound,
			expectedResponseBody:      "{\"error\":\"сайт с указанным ID не найден\"}\n",
		},
		{
			name:                      "Error storage",
			id:                        exampleUUID.String(),
			executeErr:                domain.ErrStorage,
			expectedCallCount:         1,
			expectedContentTypeHeader: "application/json",
			expectedStatusCode:        http.StatusInternalServerError,
			expectedResponseBody:      "{\"error\":\"внутренняя ошибка сервера\"}\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := mocks.NewController()
			c.Return("Execute", testCase.expectedSite, testCase.executeErr)
			getStatusUseCase := mocks.NewMockGetStatusByIDUseCase(c)

			handler := &Handler{
				getStatusUseCase: &getStatusUseCase,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/sites/"+testCase.id+"/status", nil)
			req.SetPathValue("id", testCase.id)

			handler.GetStatusByID(w, req)

			if got := getStatusUseCase.CallCount("Execute"); got != testCase.expectedCallCount {
				t.Fatalf("Expected %d, received %d", testCase.expectedCallCount, got)
			}
			if header := w.Header().Get("Content-Type"); testCase.expectedContentTypeHeader != header {
				t.Fatalf("Expected %s, received %s", testCase.expectedContentTypeHeader, header)
			}
			if testCase.expectedStatusCode != w.Code {
				t.Fatalf("Expected %d, received %d", testCase.expectedStatusCode, w.Code)
			}
			if body := w.Body.String(); testCase.expectedResponseBody != body {
				t.Fatalf("Expected %s, received %s", testCase.expectedResponseBody, body)
			}
		})
	}
}
