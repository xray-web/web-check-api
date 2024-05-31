package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
)

func TestDnsHandler(t *testing.T) {
	router := gin.Default()
	dnsCtrl := &controllers.DnsController{}
	router.GET("/dns", dnsCtrl.DnsHandler)

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
			url:          "https://example.com",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/dns?url="+tc.url, nil)
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
