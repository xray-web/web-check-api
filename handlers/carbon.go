package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CarbonData struct {
	Statistics struct {
		AdjustedBytes float64 `json:"adjustedBytes"`
		Energy        float64 `json:"energy"`
		Co2           struct {
			Grid struct {
				Grams  float64 `json:"grams"`
				Litres float64 `json:"litres"`
			} `json:"grid"`
			Renewable struct {
				Grams  float64 `json:"grams"`
				Litres float64 `json:"litres"`
			} `json:"renewable"`
		} `json:"co2"`
	} `json:"statistics"`
	CleanerThan int    `json:"cleanerThan"`
	ScanUrl     string `json:"scanUrl"`
}

// Function to get the HTML size of the website
func getHtmlSize(ctx context.Context, url string) (int, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get HTML size: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return len(body), nil
}

// Function to get the carbon data based on the HTML size
func getCarbonData(ctx context.Context, sizeInBytes int) (*CarbonData, error) {
	const carbonDataUrl = "https://api.websitecarbon.com/data"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, carbonDataUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("bytes", strconv.Itoa(sizeInBytes))
	q.Add("green", "0")
	req.URL.RawQuery = q.Encode()

	client := http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get carbon data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var carbonData CarbonData
	if err := json.Unmarshal(body, &carbonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal carbon data: %w", err)
	}

	return &carbonData, nil
}

func HandleCarbon() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		url := rawURL.String()
		sizeInBytes, err := getHtmlSize(r.Context(), url)
		if err != nil {
			JSONError(w, fmt.Errorf("error getting HTML size: %v", err), http.StatusInternalServerError)
			return
		}

		carbonData, err := getCarbonData(r.Context(), sizeInBytes)
		if err != nil {
			JSONError(w, fmt.Errorf("error getting carbon data: %v", err), http.StatusInternalServerError)
			return
		}

		if carbonData.Statistics.AdjustedBytes == 0 {
			JSON(w, KV{"skipped": "Not enough info to get carbon data"}, http.StatusOK)
			return
		}
		if carbonData.Statistics.Energy == 0 {
			JSON(w, KV{"skipped": "Not enough info to get carbon data"}, http.StatusOK)
			return
		}

		carbonData.ScanUrl = url
		JSON(w, carbonData, http.StatusOK)
	})
}
