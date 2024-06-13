package handlers

import (
	"fmt"
	"net/http"

	"github.com/xray-web/web-check-api/checks"
)

func HandleCarbon(c *checks.Carbon) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		url := rawURL.String()
		sizeInBytes, err := c.HtmlSize(r.Context(), url)
		if err != nil {
			JSONError(w, fmt.Errorf("error getting HTML size: %v", err), http.StatusInternalServerError)
			return
		}

		carbonData, err := c.CarbonData(r.Context(), sizeInBytes)
		if err != nil {
			JSONError(w, fmt.Errorf("error getting carbon data: %v", err), http.StatusInternalServerError)
			return
		}

		if carbonData.Statistics.AdjustedBytes == 0 {
			JSON(w, KV{"skipped": "Not enough info to get carbon data"}, http.StatusOK)
			return
		}
		if carbonData.Statistics.Energy == 0 {
			JSON(w, KV{"skipped": "Not enough info to get carbon data"}, http.StatusOK)
			return
		}

		carbonData.ScanUrl = url
		JSON(w, carbonData, http.StatusOK)
	})
}
