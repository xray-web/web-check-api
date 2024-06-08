package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestHandleGetSocialTags(t *testing.T) {
	// t.Parallel()
	tests := []struct {
		name           string
		urlParam       string
		mockResponse   string
		mockStatusCode int
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "missing URL parameter"},
		},
		{
			name:           "Valid URL with social tags",
			urlParam:       "http://example.com",
			mockResponse:   `<html><head><title>Example Domain</title><meta name="description" content="Example description"><meta property="og:title" content="Example OG Title"></head><body></body></html>`,
			mockStatusCode: http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
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
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			defer gock.Off()

			if tc.urlParam != "" {
				gock.New(tc.urlParam).
					Reply(tc.mockStatusCode).
					BodyString(tc.mockResponse)
			}

			req := httptest.NewRequest("GET", "/social-tags?url="+tc.urlParam, nil)
			rec := httptest.NewRecorder()
			HandleGetSocialTags().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
