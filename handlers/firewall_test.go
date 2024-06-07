package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleFirewall(t *testing.T) {
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
			req := httptest.NewRequest("GET", "/firewall?url="+tc.url, nil)
			rec := httptest.NewRecorder()

			HandleFirewall().ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, rec.Code)
			}
		})
	}
}
