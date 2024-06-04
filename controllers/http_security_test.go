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

func TestHttpSecurityHandler(t *testing.T) {
	router := gin.Default()

	ctrl := &controllers.HttpSecurityController{}

	router.GET("/check-http-security", ctrl.HttpSecurityHandler)

	req, err := http.NewRequest("GET", "/check-http-security?url=www.google.com", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.HTTPSecurityResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.False(t, response.StrictTransportPolicy)
	assert.True(t, response.XFrameOptions)
	assert.False(t, response.XContentTypeOptions)
	assert.True(t, response.XXSSProtection)
	assert.False(t, response.ContentSecurityPolicy)

}
