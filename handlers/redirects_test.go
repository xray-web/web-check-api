package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetRedirects(t *testing.T) {
	t.Parallel()

	t.Run("Missing URL parameter", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "/redirects?url=", nil)
		rec := httptest.NewRecorder()

		HandleGetRedirects().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response KV
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, KV{"error": "missing URL parameter"}, response)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "/redirects?url=invalid-url", nil)
		rec := httptest.NewRecorder()

		HandleGetRedirects().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response KV
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], KV{"error": "error: Get \"http://invalid-url\""}["error"])
	})

	t.Run("Valid URL with no redirects", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "/redirects?url=example.com", nil)
		rec := httptest.NewRecorder()

		HandleGetRedirects().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response KV
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, KV{"redirects": []interface{}{"http://example.com"}}, response)
	})
}
