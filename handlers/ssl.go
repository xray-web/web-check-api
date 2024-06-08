package handlers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func HandleSSL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			JSONError(w, errors.New("url query parameter is required"), http.StatusBadRequest)
			return
		}

		parsedURL, err := url.Parse(urlParam)
		if err != nil {
			JSONError(w, errors.New("invalid URL format"), http.StatusBadRequest)
			return
		}

		options := &tls.Config{
			ServerName:         parsedURL.Hostname(),
			InsecureSkipVerify: true, // Skip certificate validation
		}

		conn, err := tls.Dial("tcp", parsedURL.Host+":443", options)
		if err != nil {
			JSONError(w, fmt.Errorf("error establishing TLS connection: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		state := conn.ConnectionState()
		if len(state.PeerCertificates) == 0 {
			JSONError(w, errors.New("no certificate presented by the server"), http.StatusInternalServerError)
			return
		}

		cert := state.PeerCertificates[0]

		// Remove the raw field from the certificate
		cert.Raw = nil

		JSON(w, cert, http.StatusOK)
	})
}
