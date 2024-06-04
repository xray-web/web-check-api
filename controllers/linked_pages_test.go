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

func TestGetLinksHandler(t *testing.T) {
	router := gin.Default()

	ctrl := &controllers.GetLinksController{}

	router.GET("/get-links", ctrl.GetLinksHandler)

	req, err := http.NewRequest("GET", "/get-links?url=www.google.com", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.LinkResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)
	assert.NotNil(t, response.Internal)
	assert.NotNil(t, response.External)
}
