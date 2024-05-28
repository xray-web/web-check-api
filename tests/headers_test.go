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

func TestGetHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := &controllers.HeaderController{}

	router := gin.Default()
	router.GET("/headers", ctrl.GetHeaders)

	t.Run("url parameter is missing", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/headers", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "url parameter is required"}`, resp.Body.String())
	})

	t.Run("invalid url format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/headers?url=invalid-url", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("valid url", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Custom-Header", "value")
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		req, _ := http.NewRequest(http.MethodGet, "/headers?url="+mockServer.URL, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
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
