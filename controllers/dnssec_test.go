package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
)

func TestDnssecHandler(t *testing.T) {
	router := gin.Default()
	dnssecCtrl := &controllers.DnssecController{}
	router.GET("/dnssec", dnssecCtrl.DnssecHandler)

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
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/dnssec?url="+tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, resp.Code)
			}
		})
	}
}
