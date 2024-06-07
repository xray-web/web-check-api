package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleHsts(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/check-hsts?url=example.com", nil)
	rec := httptest.NewRecorder()
	HandleHsts().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response HSTSResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)
	assert.Equal(t, "Site does not serve any HSTS headers.", response.Message)
	assert.False(t, response.Compatible)
	assert.Empty(t, response.HSTSHeader)
}
