package handlers

import (
	"net/http"

	"github.com/xray-web/web-check-api/checks"
)

func HandleHsts(h *checks.Hsts) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := h.Validate(r.Context(), rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
