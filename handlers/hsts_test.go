package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleHsts(t *testing.T) {
	t.Parallel()

	t.Run("missing URL parameter", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/check-hsts", nil)
		rec := httptest.NewRecorder()

		HandleHsts(nil).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing URL parameter"}`, rec.Body.String())
	})
}
