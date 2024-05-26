package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avonture/go-wappalyzer"
	"github.com/gin-gonic/gin"
)

type TechStackController struct{}

func (ctrl *TechStackController) TechStackHandler(c *gin.Context) {
	url := c.Query("url")

	results, err := GetTechStack(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the results as JSON
	c.JSON(http.StatusOK, results)
}

func GetTechStack(url string) (string, error) {
	options := wappalyzer.Options{}

	wappalyzer, err := wappalyzer.New(options)
	if err != nil {
		return "", err
	}
	defer wappalyzer.Destroy()

	err = wappalyzer.Init()
	if err != nil {
		return "", err
	}

	headers := make(map[string]string)
	storage := wappalyzer.Storage{
		Local:   make(map[string]string),
		Session: make(map[string]string),
	}

	site, err := wappalyzer.Open(url, headers, storage)
	if err != nil {
		return "", err
	}

	results, err := site.Analyze()
	if err != nil {
		return "", err
	}

	if len(results.Technologies) == 0 {
		return "", errors.New("Unable to find any technologies for site")
	}

	// Marshal the results object to JSON
	jsonResults, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	// Return the JSON string
	return string(jsonResults), nil
}
