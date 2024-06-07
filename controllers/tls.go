package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TlsController struct{}

const MOZILLA_TLS_OBSERVATORY_API = "https://tls-observatory.services.mozilla.com/api/v1"

type ScanResponse struct {
	ScanID int `json:"scan_id"`
}

func (ctrl *TlsController) TlsHandler(c *gin.Context) {
	rawUrl := c.Query("url")
	if rawUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
		rawUrl = "http://" + rawUrl
	}

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	domain := parsedUrl.Hostname()
	scanResponse, err := initiateScan(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if scanResponse.ScanID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get scan_id from TLS Observatory"})
		return
	}

	result, err := getScanResults(scanResponse.ScanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func initiateScan(domain string) (*ScanResponse, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s/scan", MOZILLA_TLS_OBSERVATORY_API), url.Values{"target": {domain}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scanResponse ScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&scanResponse); err != nil {
		return nil, err
	}

	return &scanResponse, nil
}

func getScanResults(scanID int) (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/results?id=%d", MOZILLA_TLS_OBSERVATORY_API, scanID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func HandleTls() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}
	})
}
