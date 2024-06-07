package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleDNSServer(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Missing URL",
			url:          "",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"missing URL parameter"}`,
		},
		{
			name:         "Valid URL",
			url:          "https://example.com",
			expectedCode: http.StatusOK,
			expectedBody: `{"domain":"example.com","dns":[{"address":"93.184.215.14","hostname":null,"dohDirectSupports":false}]}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/dns?url="+tc.url, nil)
			rec := httptest.NewRecorder()

			HandleDNSServer().ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, rec.Code)
			}

			if strings.TrimSpace(rec.Body.String()) != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
