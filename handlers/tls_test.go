package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/checks"
	"github.com/xray-web/web-check-api/testutils"
)

func TestHandleTLS(t *testing.T) {
	t.Parallel()

	t.Run("Missing URL parameter", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "/tls?url=", nil)
		rec := httptest.NewRecorder()

		HandleTLS(checks.NewTls(nil)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"error": "missing URL parameter"}, responseBody)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(
			testutils.Response(http.StatusOK, []byte(`{"scan_id": 0}`)),
		)
		req := httptest.NewRequest("GET", "/tls?url=http://invalid-url", nil)
		rec := httptest.NewRecorder()

		HandleTLS(checks.NewTls(client)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"error": "failed to get scan_id from TLS Observatory"}, responseBody)
	})

}
