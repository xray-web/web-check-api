package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRedirectsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
			expectedBody:   gin.H{"error": "url parameter is required"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "invalid-url",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "error: Get \"http://invalid-url\""},
		},
		{
			name:           "Valid URL with no redirects",
			urlParam:       "example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"redirects": []interface{}{"http://example.com"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			ctrl := &controllers.RedirectsController{}
			r.GET("/redirects", ctrl.GetRedirectsHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/redirects?url="+tt.urlParam, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			if tt.name == "Invalid URL" {
				assert.Contains(t, response["error"], tt.expectedBody["error"])
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

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
			expectedBody:   controllers.KV{"error": "missing URL parameter"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "invalid-url",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   controllers.KV{"error": "error: Get \"http://invalid-url\""},
		},
		{
			name:           "Valid URL with no redirects",
			urlParam:       "example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   controllers.KV{"redirects": []interface{}{"http://example.com"}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/redirects?url="+tc.urlParam, nil)
			rec := httptest.NewRecorder()
			controllers.HandleGetRedirects().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			if tc.name == "Invalid URL" {
				assert.Contains(t, response["error"], tc.expectedBody["error"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
		})
	}
}
