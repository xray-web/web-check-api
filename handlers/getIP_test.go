package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/get-ip?url=example.com", nil)
	rec := httptest.NewRecorder()
	HandleGetIP().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	ip, ok := response["ip"].(string)
	assert.True(t, ok, "IP address not found in response")
	assert.NotEmpty(t, ip, "IP address is empty")

	family, ok := response["family"].(float64)
	assert.True(t, ok, "Family field not found in response")
	assert.Equal(t, float64(4), family, "Family field should be 4")
}
