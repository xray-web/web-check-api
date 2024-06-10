package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func HandleGetQuality() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		apiKey := os.Getenv("GOOGLE_CLOUD_API_KEY")
		if apiKey == "" {
			JSONError(w, errors.New("missing Google API. You need to set the `GOOGLE_CLOUD_API_KEY` environment variable"), http.StatusInternalServerError)
			return
		}

		encodedURL := url.QueryEscape(rawURL.String())
		endpoint := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&category=PERFORMANCE&category=ACCESSIBILITY&category=BEST_PRACTICES&category=SEO&category=PWA&strategy=mobile&key=%s", encodedURL, apiKey)

		resp, err := http.Get(endpoint)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errorResult map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&errorResult); err != nil {
				JSONError(w, err, http.StatusInternalServerError)
			} else {
				JSON(w, errorResult, resp.StatusCode)
			}
			return
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
