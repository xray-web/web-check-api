package handlers

import (
	"fmt"
	"net/http"
)

type HTTPSecurityResponse struct {
	StrictTransportPolicy bool `json:"strictTransportPolicy"`
	XFrameOptions         bool `json:"xFrameOptions"`
	XContentTypeOptions   bool `json:"xContentTypeOptions"`
	XXSSProtection        bool `json:"xXSSProtection"`
	ContentSecurityPolicy bool `json:"contentSecurityPolicy"`
}

func checkHTTPSecurity(url string) (HTTPSecurityResponse, error) {
	fullURL := "http://" + url

	resp, err := http.Get(fullURL)
	if err != nil {
		return HTTPSecurityResponse{}, fmt.Errorf("error making request: %s", err.Error())
	}
	defer resp.Body.Close()

	headers := resp.Header

	return HTTPSecurityResponse{
		StrictTransportPolicy: headers.Get("strict-transport-security") != "",
		XFrameOptions:         headers.Get("x-frame-options") != "",
		XContentTypeOptions:   headers.Get("x-content-type-options") != "",
		XXSSProtection:        headers.Get("x-xss-protection") != "",
		ContentSecurityPolicy: headers.Get("content-security-policy") != "",
	}, nil
}

func HandleHttpSecurity() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := checkHTTPSecurity(url)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
