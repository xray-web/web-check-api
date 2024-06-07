package handlers

import (
	"net/http"
)

func HandleGetHeaders() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		// Ensure the URL has a scheme
		if !(len(url) >= 7 && (url[:7] == "http://" || url[:8] == "https://")) {
			url = "http://" + url
		}

		resp, err := http.Get(url)
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
