package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBlockListsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := &controllers.BlockListsController{}
	router := gin.Default()
	router.GET("/blocklists", ctrl.BlockListsHandler)

	t.Run("missing URL parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/blocklists", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "Missing URL parameter"}`, resp.Body.String())
	})

	t.Run("blocked domain", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/blocklists?url=http://blocked.example.com", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

	})

	t.Run("unblocked domain", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/blocklists?url=http://unblocked.example.com", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

	})
}
