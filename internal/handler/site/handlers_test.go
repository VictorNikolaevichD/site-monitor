package site

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

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
