package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleTraceRoute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		urlParam       string
		expectedStatus int
		expectedBody   KV
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   KV{"error": "missing URL parameter"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/trace-route?url="+tc.urlParam, nil)
			rec := httptest.NewRecorder()
			HandleTraceRoute().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var responseBody KV
			err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
