package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLegacyRankHandler(t *testing.T) {
	router := gin.Default()

	ctrl := &controllers.LegacyRankController{}

	router.GET("/legacy-rank", ctrl.LegacyRankHandler)

	req, err := http.NewRequest("GET", "/legacy-rank?url=www.google.com", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.RankResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.Equal(t, "www.google.com", response.Domain)

	assert.True(t, response.IsFound || !response.IsFound)
}
