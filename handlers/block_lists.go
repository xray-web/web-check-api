package handlers

import (
	"net/http"

	"github.com/xray-web/web-check-api/checks"
)

func HandleBlockLists(b *checks.BlockList) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}
		list := b.BlockedServers(r.Context(), rawURL.Hostname())
		JSON(w, list, http.StatusOK)
	})
}
