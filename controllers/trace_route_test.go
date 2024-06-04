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

			ctrl := &controllers.TraceRouteController{}
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
