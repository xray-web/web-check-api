package handlers

import (
	"net/http"

	"github.com/xray-web/web-check-api/checks"
)

func HandleGetHeaders(h *checks.Headers) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		headers, err := h.List(r.Context(), rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
		}

		JSON(w, headers, http.StatusOK)
	})
}
