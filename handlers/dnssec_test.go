package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleDnsSec(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{
			name:         "Missing URL",
			url:          "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Valid URL",
			url:          "example.com",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/dnssec?url="+tc.url, nil)
			rec := httptest.NewRecorder()

			HandleDnsSec().ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, rec.Code)
			}
		})
	}
}
