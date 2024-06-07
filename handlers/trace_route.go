package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aeden/traceroute"
	"github.com/gin-gonic/gin"
)

type TraceRouteController struct{}

func (ctrl *TraceRouteController) TracerouteHandler(c *gin.Context) {
	urlString := c.Query("url")
	if urlString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "http://" + urlString
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL provided"})
		return
	}

	host := parsedURL.Hostname()
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL provided: hostname not found"})
		return
	}

	result, err := traceroute.Traceroute(host, &traceroute.TracerouteOptions{})
	if err != nil {
		errorMessage := fmt.Sprintf("Error performing traceroute: %s", err)
		fmt.Println(errorMessage) // Log the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error performing traceroute"})
		return
	}

	var response []string
	for _, hop := range result.Hops {
		response = append(response, fmt.Sprintf("%d. %s", hop.TTL, hop.Address))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Traceroute completed!", "hops": response})
}

func HandleTraceRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlString := r.URL.Query().Get("url")
		if urlString == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
			urlString = "http://" + urlString
		}

		parsedURL, err := url.Parse(urlString)
		if err != nil {
			JSONError(w, ErrInvalidURL, http.StatusBadRequest)
			return
		}

		host := parsedURL.Hostname()
		if host == "" {
			JSONError(w, errors.New("invalid URL provided: hostname not found"), http.StatusBadRequest)
			return
		}

		result, err := traceroute.Traceroute(host, &traceroute.TracerouteOptions{})
		if err != nil {
			JSONError(w, errors.New("error performing traceroute"), http.StatusInternalServerError)
			return
		}

		var response []string
		for _, hop := range result.Hops {
			response = append(response, fmt.Sprintf("%d. %s", hop.TTL, hop.Address))
		}

		JSON(w, KV{"message": "Traceroute completed!", "hops": response}, http.StatusOK)
	})
}
