package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetRank(t *testing.T) {
	// t.Parallel()
	tests := []struct {
		name           string
		urlParam       string
		mockResponse   interface{}
		mockStatusCode int
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "missing URL parameter"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "invalid-url",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "unable to fetch rank, Get \"https://tranco-list.eu/api/ranks/domain/invalid-url\": no responder found"},
		},
		{
			name:           "Valid request with rank found",
			urlParam:       "http://example.com",
			mockResponse:   map[string]interface{}{"ranks": []interface{}{map[string]interface{}{"rank": 1.0}}},
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"ranks": []interface{}{map[string]interface{}{"rank": 1.0}}},
		},
		{
			name:           "Valid request with no rank found",
			urlParam:       "http://example.com",
			mockResponse:   map[string]interface{}{"ranks": []interface{}{}},
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"skipped": "Skipping, as example.com isn't ranked in the top 100 million sites yet."},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("TRANCO_API_KEY", "test_api_key")
			os.Setenv("TRANCO_USERNAME", "test_username")

			if tc.mockResponse != nil {
				httpmock.RegisterResponder("GET", "https://tranco-list.eu/api/ranks/domain/example.com",
					func(req *http.Request) (*http.Response, error) {
						res, err := httpmock.NewJsonResponse(tc.mockStatusCode, tc.mockResponse)
						if err != nil {
							return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
						}
						return res, nil
					})
			}

			req, _ := http.NewRequest("GET", "/rank?url="+tc.urlParam, nil)
			w := httptest.NewRecorder()

			HandleGetRank().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
