package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeaderController struct{}

func (ctrl *HeaderController) GetHeaders(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Ensure the URL has a scheme
	if !(len(url) >= 7 && (url[:7] == "http://" || url[:8] == "https://")) {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Copying headers from the response
	headers := make(map[string]interface{})
	for key, values := range resp.Header {
		if len(values) > 1 {
			headers[key] = values
		} else {
			headers[key] = values[0]
		}
	}

	c.JSON(http.StatusOK, headers)
}
