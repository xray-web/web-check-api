package controllers

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetIPController struct{}

func lookupAsync(address string) (map[string]interface{}, error) {
	ip, err := net.LookupIP(address)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if len(ip) > 0 {
		result["ip"] = ip[0].String()
		result["family"] = 4
	} else {
		result["ip"] = ""
		result["family"] = nil
	}

	return result, nil
}

func (ctrl *GetIPController) GetIPHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	address := strings.ReplaceAll(strings.ReplaceAll(url, "https://", ""), "http://", "")
	result, err := lookupAsync(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func HandleGetIP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		address := strings.ReplaceAll(strings.ReplaceAll(url, "https://", ""), "http://", "")
		result, err := lookupAsync(address)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
