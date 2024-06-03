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

func TestGetPortsHandler(t *testing.T) {
	router := gin.Default()
	ctrl := &controllers.PortsController{}
	router.GET("/check-ports", ctrl.GetPortsHandler)

	tests := []struct {
		url          string
		expectedCode int
		expectedBody map[string][]int
	}{
		{
			url:          "open.com",
			expectedCode: http.StatusOK,
		},
		{
			url:          "closed.com",
			expectedCode: http.StatusOK,
		},
		{
			url:          "mixed.com",
			expectedCode: http.StatusOK,
		},
		{
			url:          "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/check-ports?url="+tt.url, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string][]int
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				openPorts, ok1 := response["openPorts"]
				failedPorts, ok2 := response["failedPorts"]

				assert.True(t, ok1)
				assert.True(t, ok2)
				assert.NotNil(t, openPorts)
				assert.NotNil(t, failedPorts)
			} else {
				var response map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "url parameter is required")
			}
		})
	}
}
