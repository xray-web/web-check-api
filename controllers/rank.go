package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RankController struct{}

func (ctrl *RankController) GetRankHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url // Assuming HTTP by default
	}

	domain, err := getDomain(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auth := getAuth()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://tranco-list.eu/api/ranks/domain/%s", domain), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if auth != nil {
		req.SetBasicAuth(auth["username"], auth["password"])
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to fetch rank, %v", err)})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ranks, ok := result["ranks"].([]interface{}); !ok || len(ranks) == 0 {
		c.JSON(http.StatusOK, gin.H{"skipped": fmt.Sprintf("Skipping, as %s isn't ranked in the top 100 million sites yet.", domain)})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getDomain(rawurl string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawurl)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL")
	}
	return parsedURL.Hostname(), nil
}

func getAuth() map[string]string {
	apiKey := os.Getenv("TRANCO_API_KEY")
	username := os.Getenv("TRANCO_USERNAME")

	if apiKey != "" && username != "" {
		return map[string]string{
			"username": username,
			"password": apiKey,
		}
	}

	return nil
}
