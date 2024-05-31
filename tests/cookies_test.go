package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
)

func TestCookiesHandler(t *testing.T) {
	router := gin.Default()
	cookiesCtrl := &controllers.CookiesController{}
	router.GET("/cookies", cookiesCtrl.CookiesHandler)

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
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/cookies?url="+tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, resp.Code)
			}

			if tc.expectedError != "" && !strings.Contains(resp.Body.String(), tc.expectedError) {
				t.Errorf("Expected error message '%s' not found in response body", tc.expectedError)
			}
		})
	}
}
