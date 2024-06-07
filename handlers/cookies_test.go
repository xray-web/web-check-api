package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerCookies(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		url           string
		expectedCode  int
		expectedError string
	}{
		{
			name:         "Missing URL",
			url:          "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid URL",
			url:          "invalid_url",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/cookies?url="+tc.url, nil)
			rec := httptest.NewRecorder()
			HandleCookies().ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, rec.Code)
			}

			if tc.expectedError != "" && !strings.Contains(rec.Body.String(), tc.expectedError) {
				t.Errorf("Expected error message '%s' not found in response body", tc.expectedError)
			}
		})
	}
}
