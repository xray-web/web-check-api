package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
