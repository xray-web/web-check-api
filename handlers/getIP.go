package handlers

import (
	"net"
	"net/http"
)

func lookupAsync(address string) (map[string]interface{}, error) {
	ip, err := net.LookupIP(address)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if len(ip) > 0 {
		result["ip"] = ip[0].String()
		result["family"] = 4
	} else {
		result["ip"] = ""
		result["family"] = nil
	}

	return result, nil
}

func HandleGetIP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := lookupAsync(rawURL.Hostname())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
