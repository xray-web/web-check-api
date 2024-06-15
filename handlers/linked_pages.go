package handlers

import (
	"fmt"
	"net/http"

	"github.com/xray-web/web-check-api/checks"
)

func HandleGetLinks(l *checks.LinkedPages) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		links, err := l.GetLinkedPages(r.Context(), rawURL)
		if err != nil {
			JSONError(w, fmt.Errorf("error getting linked pages: %v", err), http.StatusInternalServerError)
			return
		}

		if len(links.Internal) == 0 && len(links.External) == 0 {
			JSON(w, KV{
				"skipped": `No internal or external links found. 
				This may be due to the website being dynamically rendered, using a client-side framework (like React), and without SSR enabled. 
				That would mean that the static HTML returned from the HTTP request doesn't contain any meaningful content for Web-Check to analyze. 
				You can rectify this by using a headless browser to render the page instead.`,
			}, http.StatusOK)
			return
		}

		JSON(w, links, http.StatusOK)
	})
}
