package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDnsServerHandler(t *testing.T) {
	router := gin.Default()
	dnsCtrl := &DnsServerController{}
	router.GET("/dns", dnsCtrl.DnsServerHandler)

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
			expectedBody: `{"error":"url parameter is required"}`,
		},
		{
			name:         "Valid URL",
			url:          "https://example.com",
			expectedCode: http.StatusOK,
			expectedBody: `{"domain":"example.com","dns":[{"address":"93.184.215.14","hostname":null,"dohDirectSupports":false}]}`,
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

			if strings.TrimSpace(resp.Body.String()) != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, strings.TrimSpace(resp.Body.String()))
			}
		})
	}
}

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

			if rec.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, rec.Code)
			}

			if strings.TrimSpace(rec.Body.String()) != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
