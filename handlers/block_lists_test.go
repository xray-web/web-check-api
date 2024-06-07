package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBlockListsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := &BlockListsController{}
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

func TestHandleBlockLists(t *testing.T) {
	t.Parallel()

	t.Run("missing URL parameter", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/blocklists", nil)
		rec := httptest.NewRecorder()

		HandleBlockLists().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing URL parameter"}`, rec.Body.String())
	})

	t.Run("blocked domain", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/blocklists?url=http://blocked.example.com", nil)
		rec := httptest.NewRecorder()

		HandleBlockLists().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

	})

	t.Run("unblocked domain", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/blocklists?url=http://unblocked.example.com", nil)
		rec := httptest.NewRecorder()

		HandleBlockLists().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

	})
}
