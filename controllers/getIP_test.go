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

func TestGetIPHandler(t *testing.T) {
	router := gin.Default()
	ctrl := &controllers.GetIPController{}
	router.GET("/get-ip", ctrl.GetIPHandler)

	req, err := http.NewRequest("GET", "/get-ip?url=example.com", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	ip, ok := response["ip"].(string)
	assert.True(t, ok, "IP address not found in response")
	assert.NotEmpty(t, ip, "IP address is empty")

	family, ok := response["family"].(float64)
	assert.True(t, ok, "Family field not found in response")
	assert.Equal(t, float64(4), family, "Family field should be 4")
}

func TestHandleGetIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/get-ip?url=example.com", nil)
	rec := httptest.NewRecorder()
	controllers.HandleGetIP().ServeHTTP(rec, req)

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
