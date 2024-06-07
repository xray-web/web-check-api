package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTraceRouteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		urlParam       string
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "url parameter is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			ctrl := &TraceRouteController{}
			router.GET("/trace-route", ctrl.TracerouteHandler)

			req, _ := http.NewRequest("GET", "/trace-route?url="+tt.urlParam, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody gin.H
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestHandleTraceRoute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		urlParam       string
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "missing URL parameter"},
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

			var responseBody gin.H
			err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
