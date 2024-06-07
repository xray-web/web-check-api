package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func getRedirects(rawurl string) ([]string, error) {
	redirects := []string{rawurl}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 12 {
				return fmt.Errorf("stopped after 12 redirects")
			}
			redirects = append(redirects, req.URL.String())
			return nil
		},
	}

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer resp.Body.Close()

	return redirects, nil
}

func HandleGetRedirects() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlParam, "http://") && !strings.HasPrefix(urlParam, "https://") {
			urlParam = "http://" + urlParam // Assuming HTTP by default
		}

		redirects, err := getRedirects(urlParam)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, KV{"redirects": redirects}, http.StatusOK)
	})
}
