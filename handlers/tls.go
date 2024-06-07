package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const MOZILLA_TLS_OBSERVATORY_API = "https://tls-observatory.services.mozilla.com/api/v1"

type ScanResponse struct {
	ScanID int `json:"scan_id"`
}

func initiateScan(domain string) (*ScanResponse, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s/scan", MOZILLA_TLS_OBSERVATORY_API), url.Values{"target": {domain}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scanResponse ScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&scanResponse); err != nil {
		return nil, err
	}

	return &scanResponse, nil
}

func getScanResults(scanID int) (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/results?id=%d", MOZILLA_TLS_OBSERVATORY_API, scanID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func HandleTLS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUrl := r.URL.Query().Get("url")
		if rawUrl == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
			rawUrl = "http://" + rawUrl
		}

		parsedUrl, err := url.Parse(rawUrl)
		if err != nil {
			JSONError(w, ErrInvalidURL, http.StatusBadRequest)
			return
		}

		domain := parsedUrl.Hostname()
		scanResponse, err := initiateScan(domain)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if scanResponse.ScanID == 0 {
			JSONError(w, errors.New("failed to get scan_id from TLS Observatory"), http.StatusInternalServerError)
			return
		}

		result, err := getScanResults(scanResponse.ScanID)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
