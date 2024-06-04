package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRankHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := &controllers.RankController{}
	router.GET("/rank", ctrl.GetRankHandler)

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
			expectedBody:   map[string]interface{}{"error": "url parameter is required"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "invalid-url",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "Unable to fetch rank, Get \"https://tranco-list.eu/api/ranks/domain/invalid-url\": no responder found"},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TRANCO_API_KEY", "test_api_key")
			os.Setenv("TRANCO_USERNAME", "test_username")

			if tt.mockResponse != nil {
				httpmock.RegisterResponder("GET", "https://tranco-list.eu/api/ranks/domain/example.com",
					func(req *http.Request) (*http.Response, error) {
						res, err := httpmock.NewJsonResponse(tt.mockStatusCode, tt.mockResponse)
						if err != nil {
							return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
						}
						return res, nil
					})
			}

			req, _ := http.NewRequest("GET", "/rank?url="+tt.urlParam, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
