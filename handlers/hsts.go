package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type HSTSResponse struct {
	Message    string `json:"message"`
	Compatible bool   `json:"compatible"`
	HSTSHeader string `json:"hstsHeader"`
}

func checkHSTS(url string) (HSTSResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return HSTSResponse{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return HSTSResponse{}, fmt.Errorf("error making request: %s", err.Error())
	}
	defer resp.Body.Close()

	hstsHeader := resp.Header.Get("strict-transport-security")
	if hstsHeader == "" {
		return HSTSResponse{Message: "Site does not serve any HSTS headers."}, nil
	}

	maxAgeMatch := regexp.MustCompile(`max-age=(\d+)`).FindStringSubmatch(hstsHeader)
	if maxAgeMatch == nil || len(maxAgeMatch) < 2 || maxAgeMatch[1] == "" || maxAgeMatch[1] < "10886400" {
		return HSTSResponse{Message: "HSTS max-age is less than 10886400."}, nil
	}

	if !strings.Contains(hstsHeader, "includeSubDomains") {
		return HSTSResponse{Message: "HSTS header does not include all subdomains."}, nil
	}

	if !strings.Contains(hstsHeader, "preload") {
		return HSTSResponse{Message: "HSTS header does not contain the preload directive."}, nil
	}

	return HSTSResponse{Message: "Site is compatible with the HSTS preload list!", Compatible: true, HSTSHeader: hstsHeader}, nil
}

func HandleHsts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := checkHSTS(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
