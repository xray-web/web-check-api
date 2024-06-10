package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

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

func HandleGetRank() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		domain, err := getDomain(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}

		auth := getAuth()

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("https://tranco-list.eu/api/ranks/domain/%s", domain), nil)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if auth != nil {
			req.SetBasicAuth(auth["username"], auth["password"])
		}

		resp, err := client.Do(req)
		if err != nil {
			JSONError(w, fmt.Errorf("unable to fetch rank, %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if ranks, ok := result["ranks"].([]interface{}); !ok || len(ranks) == 0 {
			JSON(w, KV{"skipped": fmt.Sprintf("Skipping, as %s isn't ranked in the top 100 million sites yet.", domain)}, http.StatusOK)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
