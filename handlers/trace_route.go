package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aeden/traceroute"
)

func HandleTraceRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := traceroute.Traceroute(rawURL.Host, &traceroute.TracerouteOptions{})
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
