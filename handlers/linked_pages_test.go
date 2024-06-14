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

func TestHandleGetLinks(t *testing.T) {
	t.Parallel()

	t.Run("missing URL parameter", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "/legacy-rank?url=", nil)
		rec := httptest.NewRecorder()

		HandleGetLinks(nil).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error": "missing URL parameter"}`, rec.Body.String())
	})

	t.Run("failed to get page body", func(t *testing.T) {
		t.Parallel()
		client := testutils.MockClient(testutils.Response(http.StatusInternalServerError, nil))

		req := httptest.NewRequest("GET", "/legacy-rank?url=http://test.com", nil)
		rec := httptest.NewRecorder()
		HandleGetLinks(checks.NewLinkedPages(client)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response KV
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "error getting linked pages")
	})

	t.Run("no linked pages", func(t *testing.T) {
		t.Parallel()

		testHTML := []byte(`<html><body>Test</body></html>`)
		client := testutils.MockClient(testutils.Response(http.StatusOK, testHTML))
		req := httptest.NewRequest("GET", "/legacy-rank?url=http://test.com", nil)
		rec := httptest.NewRecorder()

		HandleGetLinks(checks.NewLinkedPages(client)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response KV
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["skipped"], "No internal or external links found")
	})

	t.Run("linked pages successfully found", func(t *testing.T) {
		t.Parallel()

		testHTML := []byte(`
		<a href="http://test.com/"></a>
		<a href="http://external.com/"></a>`)
		client := testutils.MockClient(testutils.Response(http.StatusOK, testHTML))
		req := httptest.NewRequest("GET", "/legacy-rank?url=http://test.com", nil)
		rec := httptest.NewRecorder()

		HandleGetLinks(checks.NewLinkedPages(client)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response checks.LinkedPagesData
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response.Internal)
		assert.NotNil(t, response.External)
	})
}
