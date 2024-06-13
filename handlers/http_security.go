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
	// fullURL := "http://" + url
	// TODO(Lissy93): does this test require we set scheme to http?
	resp, err := http.Get(url)
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
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := checkHTTPSecurity(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
