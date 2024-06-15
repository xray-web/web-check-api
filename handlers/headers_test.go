package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/checks"
)

func TestHandleGetHeaders(t *testing.T) {
	t.Parallel()

	t.Run("url parameter is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/headers", nil)
		rec := httptest.NewRecorder()
		HandleGetHeaders(nil).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing URL parameter"}`, rec.Body.String())
	})

	t.Run("invalid url format", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Custom-Header", "value")
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		req := httptest.NewRequest(http.MethodGet, "/headers?url=invalid-url", nil)
		rec := httptest.NewRecorder()
		HandleGetHeaders(checks.NewHeaders(mockServer.Client())).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
