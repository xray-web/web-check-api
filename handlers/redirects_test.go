package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetRedirects(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		urlParam       string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   KV{"error": "missing URL parameter"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "invalid-url",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   KV{"error": "error: Get \"http://invalid-url\""},
		},
		{
			name:           "Valid URL with no redirects",
			urlParam:       "example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   KV{"redirects": []interface{}{"http://example.com"}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/redirects?url="+tc.urlParam, nil)
			rec := httptest.NewRecorder()
			HandleGetRedirects().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			// TODO: break this out of table drive tests, should not use name as part of logic
			if tc.name == "Invalid URL" {
				assert.Contains(t, response["error"], tc.expectedBody["error"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
		})
	}
}
