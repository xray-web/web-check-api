package handlers

import (
	"net/http"
	"net/url"
)

func HandleGetHeaders() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")
		if rawURL == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		// Parse and validate the URL
		parsedURL, err := url.ParseRequestURI(rawURL)
		if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
			JSONError(w, ErrInvalidURL, http.StatusBadRequest)
			return
		}

		resp, err := http.Get(parsedURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
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

		JSON(w, headers, http.StatusOK)
	})
}
