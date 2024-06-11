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

func TestHandleGetSocialTags(t *testing.T) {
	t.Parallel()

	t.Run("Missing URL parameter", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "/social-tag?url=", nil)
		rec := httptest.NewRecorder()

		HandleGetSocialTags(checks.NewSocialTags(nil)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response KV
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, KV{"error": "missing URL parameter"}, response)
	})

	t.Run("Valid URL with social tags", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "/social-tags?url=example.com", nil)
		rec := httptest.NewRecorder()

		HandleGetSocialTags(checks.NewSocialTags(testutils.MockClient(testutils.Response(http.StatusOK, []byte(`<html><head><title>Example Domain</title><meta name="description" content="Example description"><meta property="og:title" content="Example OG Title"></head><body></body></html>`))))).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var responseBody KV
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, KV{
			"title":              "Example Domain",
			"description":        "Example description",
			"keywords":           "",
			"canonicalUrl":       "",
			"ogTitle":            "Example OG Title",
			"ogType":             "",
			"ogImage":            "",
			"ogUrl":              "",
			"ogDescription":      "",
			"ogSiteName":         "",
			"twitterCard":        "",
			"twitterSite":        "",
			"twitterCreator":     "",
			"twitterTitle":       "",
			"twitterDescription": "",
			"twitterImage":       "",
			"themeColor":         "",
			"robots":             "",
			"googlebot":          "",
			"generator":          "",
			"viewport":           "",
			"author":             "",
			"publisher":          "",
			"favicon":            "",
		}, responseBody)
	})

}
