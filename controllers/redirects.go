package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RedirectsController struct{}

func (ctrl *RedirectsController) GetRedirectsHandler(c *gin.Context) {
	urlParam := c.Query("url")
	if urlParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	if !strings.HasPrefix(urlParam, "http://") && !strings.HasPrefix(urlParam, "https://") {
		urlParam = "http://" + urlParam // Assuming HTTP by default
	}

	redirects, err := getRedirects(urlParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirects": redirects})
}

func getRedirects(rawurl string) ([]string, error) {
	redirects := []string{rawurl}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 12 {
				return fmt.Errorf("stopped after 12 redirects")
			}
			redirects = append(redirects, req.URL.String())
			return nil
		},
	}

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer resp.Body.Close()

	return redirects, nil
}
