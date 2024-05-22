package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type CookiesController struct{}

// Function to get cookies using chromedp
func getChromedpCookies(url string) ([]map[string]interface{}, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var cookiesStr string

	// Create a timeout context for chromedp actions
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // Increased timeout
	defer cancel()

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(`document.cookie`, &cookiesStr),
	)

	if err != nil {
		return nil, err
	}

	var cookies []map[string]interface{}
	err = json.Unmarshal([]byte(cookiesStr), &cookies)
	if err != nil {
		return nil, err
	}

	return cookies, nil
}

// Function to handle the cookies endpoint
func (ctrl *CookiesController) CookiesHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		return
	}

	// Ensure the URL includes a scheme
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	var headerCookies []string
	var clientCookies []map[string]interface{}

	// Fetch headers using http.Get
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Request failed: %v", err)})
		return
	}
	defer resp.Body.Close()

	headerCookies = resp.Header["Set-Cookie"]

	// Fetch client cookies using chromedp
	clientCookies, err = getChromedpCookies(url)
	if err != nil {
		clientCookies = nil
	}

	if len(headerCookies) == 0 && (len(clientCookies) == 0) {
		c.JSON(http.StatusOK, gin.H{"skipped": "No cookies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"headerCookies": headerCookies,
		"clientCookies": clientCookies,
	})
}
