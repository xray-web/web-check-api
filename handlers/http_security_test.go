package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
