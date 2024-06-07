package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetHeaders(t *testing.T) {
	t.Parallel()

	t.Run("url parameter is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/headers", nil)
		rec := httptest.NewRecorder()
		HandleGetHeaders().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing URL parameter"}`, rec.Body.String())
	})

	t.Run("invalid url format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/headers?url=invalid-url", nil)
		rec := httptest.NewRecorder()
		HandleGetHeaders().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "invalid URL"}`, rec.Body.String())
	})

	t.Run("valid url", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Custom-Header", "value")
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		req := httptest.NewRequest(http.MethodGet, "/headers?url="+mockServer.URL, nil)
		rec := httptest.NewRecorder()
		HandleGetHeaders().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedHeaders := map[string]interface{}{
			"Content-Type":    "application/json",
			"X-Custom-Header": "value",
		}

		for key, expectedValue := range expectedHeaders {
			assert.Equal(t, expectedValue, responseBody[key])
		}
	})
}
