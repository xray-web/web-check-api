package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptrace"
	"time"
)

type responseData struct {
	IsUp          bool    `json:"isUp"`
	DNSLookupTime float64 `json:"dnsLookupTime"`
	ResponseTime  float64 `json:"responseTime"`
	ResponseCode  int     `json:"responseCode"`
}

func fetchURL(url string) (*responseData, error) {
	if url == "" {
		return nil, errors.New("you must provide a URL query parameter!")
	}

	var dnsStart, dnsEnd, startTime time.Time
	var responseCode int

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			dnsEnd = time.Now()
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	startTime = time.Now()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Follow redirects
	for resp.StatusCode >= 300 && resp.StatusCode < 400 {
		loc, err := resp.Location()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("GET", loc.String(), nil)
		if err != nil {
			return nil, err
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	responseCode = resp.StatusCode
	if responseCode != 200 {
		return nil, errors.New("received non-success response code: " + http.StatusText(responseCode))
	}

	dnsLookupTime := dnsEnd.Sub(dnsStart).Seconds() * 1000 // Convert to milliseconds
	responseTime := time.Since(startTime).Seconds() * 1000 // Convert to milliseconds

	return &responseData{
		IsUp:          true,
		DNSLookupTime: dnsLookupTime,
		ResponseTime:  responseTime,
		ResponseCode:  responseCode,
	}, nil
}

func HandleStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		data, err := fetchURL(rawURL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
