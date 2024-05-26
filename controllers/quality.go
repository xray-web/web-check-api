package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type QualityController struct{}

func (ctrl *QualityController) GetQualityHandler(c *gin.Context) {
	urlParam := c.Query("url")
	if urlParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	apiKey := os.Getenv("GOOGLE_CLOUD_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing Google API. You need to set the `GOOGLE_CLOUD_API_KEY` environment variable"})
		return
	}

	endpoint := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&category=PERFORMANCE&category=ACCESSIBILITY&category=BEST_PRACTICES&category=SEO&category=PWA&strategy=mobile&key=%s", url.QueryEscape(urlParam), apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch the Pagespeed data"})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
