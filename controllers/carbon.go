package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CarbonController struct{}

// Function to get the HTML size of the website
func getHtmlSize(url string) (int, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get HTML size: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return len(body), nil
}

// Function to get the carbon data based on the HTML size
func getCarbonData(sizeInBytes int) (map[string]interface{}, error) {
	apiUrl := fmt.Sprintf("https://api.websitecarbon.com/data?bytes=%d&green=0", sizeInBytes)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get carbon data: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var carbonData map[string]interface{}
	if err := json.Unmarshal(body, &carbonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal carbon data: %w", err)
	}

	return carbonData, nil
}

// Handler for the /carbon endpoint
func (ctrl *CarbonController) CarbonHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		return
	}

	sizeInBytes, err := getHtmlSize(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting HTML size: %v", err)})
		return
	}

	carbonData, err := getCarbonData(sizeInBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting carbon data: %v", err)})
		return
	}

	if stats, ok := carbonData["statistics"].(map[string]interface{}); ok {
		if adjustedBytes, ok := stats["adjustedBytes"].(float64); ok && adjustedBytes == 0 {
			c.JSON(http.StatusOK, gin.H{"skipped": "Not enough info to get carbon data"})
			return
		}
		if energy, ok := stats["energy"].(float64); ok && energy == 0 {
			c.JSON(http.StatusOK, gin.H{"skipped": "Not enough info to get carbon data"})
			return
		}
	}

	carbonData["scanUrl"] = url
	c.JSON(http.StatusOK, carbonData)
}
