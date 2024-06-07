package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aeden/traceroute"
)

func HandleTraceRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlString := r.URL.Query().Get("url")
		if urlString == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
			urlString = "http://" + urlString
		}

		parsedURL, err := url.Parse(urlString)
		if err != nil {
			JSONError(w, ErrInvalidURL, http.StatusBadRequest)
			return
		}

		host := parsedURL.Hostname()
		if host == "" {
			JSONError(w, errors.New("invalid URL provided: hostname not found"), http.StatusBadRequest)
			return
		}

		result, err := traceroute.Traceroute(host, &traceroute.TracerouteOptions{})
		if err != nil {
			JSONError(w, errors.New("error performing traceroute"), http.StatusInternalServerError)
			return
		}

		var response []string
		for _, hop := range result.Hops {
			response = append(response, fmt.Sprintf("%d. %s", hop.TTL, hop.Address))
		}

		JSON(w, KV{"message": "Traceroute completed!", "hops": response}, http.StatusOK)
	})
}
