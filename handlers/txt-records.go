package handlers

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func HandleTXTRecords() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			JSONError(w, errors.New("url query parameter is required"), http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlParam, "http://") && !strings.HasPrefix(urlParam, "https://") {
			urlParam = "https://" + urlParam
		}

		parsedURL, err := url.Parse(urlParam)
		if err != nil {
			JSONError(w, errors.New("invalid URL format"), http.StatusBadRequest)
			return
		}

		txtRecords, err := resolveTXTRecords(parsedURL.Hostname())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		readableTxtRecords := parseTXTRecords(txtRecords)

		JSON(w, readableTxtRecords, http.StatusOK)
	})
}

func resolveTXTRecords(hostname string) ([]string, error) {
	txtRecords, err := net.LookupTXT(hostname)
	if err != nil {
		return nil, err
	}
	return txtRecords, nil
}

func parseTXTRecords(txtRecords []string) map[string]string {
	readableTxtRecords := make(map[string]string)
	for _, recordString := range txtRecords {
		splitRecord := strings.SplitN(recordString, "=", 2)
		if len(splitRecord) == 2 {
			key := splitRecord[0]
			value := splitRecord[1]
			readableTxtRecords[key] = value
		}
	}
	return readableTxtRecords
}
