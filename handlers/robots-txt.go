package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ParseRobotsTxt(content string) map[string][]map[string]string {
	lines := strings.Split(content, "\n")
	rules := []map[string]string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		match := ""

		if strings.HasPrefix(strings.ToLower(line), "allow:") {
			match = "Allow"
		} else if strings.HasPrefix(strings.ToLower(line), "disallow:") {
			match = "Disallow"
		} else if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			match = "User-agent"
		}

		if match != "" {
			val := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			rule := map[string]string{
				"lbl": match,
				"val": val,
			}
			rules = append(rules, rule)
		}
	}

	return map[string][]map[string]string{"robots": rules}
}

func HandleRobotsTxt() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			JSONError(w, errors.New("url query parameter is required"), http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlParam, "http://") && !strings.HasPrefix(urlParam, "https://") {
			urlParam = "http://" + urlParam
		}

		parsedURL, err := url.Parse(urlParam)
		if err != nil {
			JSONError(w, errors.New("Invalid url query parameter"), http.StatusBadRequest)
			return
		}

		robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)

		resp, err := http.Get(robotsURL)
		if err != nil {
			JSONError(w, fmt.Errorf("Error fetching robots.txt: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			JSON(w, map[string]interface{}{
				"error":      "Failed to fetch robots.txt",
				"statusCode": resp.StatusCode,
			}, resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			JSONError(w, fmt.Errorf("Error reading robots.txt: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		parsedData := ParseRobotsTxt(string(body))
		if robots, ok := parsedData["robots"]; !ok || len(robots) == 0 {
			JSON(w, map[string]string{"skipped": "No robots.txt file present, unable to continue"}, http.StatusOK)
			return
		}

		JSON(w, parsedData, http.StatusOK)
	})
}
