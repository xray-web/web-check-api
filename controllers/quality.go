package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type QualityController struct{}

func (ctrl *QualityController) GetQualityHandler(c *gin.Context) {
	urlParam := c.Query("url")
	if urlParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	formattedURL, err := formatURL(urlParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := os.Getenv("GOOGLE_CLOUD_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing Google API. You need to set the `GOOGLE_CLOUD_API_KEY` environment variable"})
		return
	}

	encodedURL := url.QueryEscape(formattedURL)
	endpoint := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&category=PERFORMANCE&category=ACCESSIBILITY&category=BEST_PRACTICES&category=SEO&category=PWA&strategy=mobile&key=%s", encodedURL, apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResult map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResult); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(resp.StatusCode, errorResult)
		}
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func formatURL(input string) (string, error) {
	// Add http scheme if missing
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}

	// Parse the URL to ensure it's valid
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Rebuild the URL to ensure it matches the required format
	return parsedURL.String(), nil
}

func HandleGetQuality() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}
	})
}
