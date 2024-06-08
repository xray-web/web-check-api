package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestHandleTLS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		urlParam         string
		mockScanResp     string
		mockScanStatus   int
		mockResultResp   string
		mockResultStatus int
		expectedStatus   int
		expectedBody     map[string]interface{}
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "missing URL parameter"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "http://invalid-url",
			mockScanResp:   `{"scan_id": 0}`,
			mockScanStatus: http.StatusOK,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to get scan_id from TLS Observatory"},
		},
		{
			name:             "Valid URL with successful scan",
			urlParam:         "http://example.com",
			mockScanResp:     `{"scan_id": 12345}`,
			mockScanStatus:   http.StatusOK,
			mockResultResp:   `{"grade": "A+"}`,
			mockResultStatus: http.StatusOK,
			expectedStatus:   http.StatusOK,
			expectedBody:     map[string]interface{}{"grade": "A+"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			defer gock.Off()

			if tc.urlParam != "" {
				gock.New(MOZILLA_TLS_OBSERVATORY_API).
					Post("/scan").
					Reply(tc.mockScanStatus).
					BodyString(tc.mockScanResp)

				if tc.mockScanStatus == http.StatusOK && tc.mockResultResp != "" {
					gock.New(MOZILLA_TLS_OBSERVATORY_API).
						Get("/results").
						MatchParam("id", "12345").
						Reply(tc.mockResultStatus).
						BodyString(tc.mockResultResp)
				}
			}

			req := httptest.NewRequest("GET", "/tls?url="+tc.urlParam, nil)
			rec := httptest.NewRecorder()
			HandleTLS().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
