package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHttpSecurityHandler(t *testing.T) {
	router := gin.Default()

	ctrl := &HttpSecurityController{}

	router.GET("/check-http-security", ctrl.HttpSecurityHandler)

	req, err := http.NewRequest("GET", "/check-http-security?url=www.google.com", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HTTPSecurityResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.False(t, response.StrictTransportPolicy)
	assert.True(t, response.XFrameOptions)
	assert.False(t, response.XContentTypeOptions)
	assert.True(t, response.XXSSProtection)
	assert.False(t, response.ContentSecurityPolicy)

}

func TestHandleHttpSecurity(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/check-http-security?url=www.google.com", nil)
	rec := httptest.NewRecorder()
	HandleHttpSecurity().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response HTTPSecurityResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.False(t, response.StrictTransportPolicy)
	assert.True(t, response.XFrameOptions)
	assert.False(t, response.XContentTypeOptions)
	assert.True(t, response.XXSSProtection)
	assert.False(t, response.ContentSecurityPolicy)
}
