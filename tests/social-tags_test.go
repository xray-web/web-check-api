package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetSocialTagsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
			expectedBody:   map[string]interface{}{"error": "url parameter is required"},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			if tt.urlParam != "" {
				gock.New(tt.urlParam).
					Reply(tt.mockStatusCode).
					BodyString(tt.mockResponse)
			}

			router := gin.Default()
			router.GET("/social-tags", func(c *gin.Context) {
				ctrl := &controllers.SocialTagsController{}
				ctrl.GetSocialTagsHandler(c)
			})

			req, _ := http.NewRequest("GET", "/social-tags?url="+tt.urlParam, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
