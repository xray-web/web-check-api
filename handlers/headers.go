package handlers

import (
	"net/http"
)

func HandleGetHeaders() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		resp, err := http.Get(rawURL.String())
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
