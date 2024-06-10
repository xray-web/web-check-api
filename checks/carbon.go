package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

type Carbon struct {
	client *http.Client
}

func NewCarbon(client *http.Client) *Carbon {
	return &Carbon{client: client}
}

// HtmlSize gets the HTML size of the website
func (c *Carbon) HtmlSize(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get HTML size: %w", err)
	}
	defer resp.Body.Close()
	// short cut to avoid reading body into RAM
	if resp.ContentLength != -1 {
		return int(resp.ContentLength), nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	return len(body), nil
}

// CarbonData gets the carbon data based on the HTML size
func (c *Carbon) CarbonData(ctx context.Context, sizeInBytes int) (*CarbonData, error) {
	const carbonDataUrl = "https://api.websitecarbon.com/data"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, carbonDataUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("bytes", strconv.Itoa(sizeInBytes))
	q.Add("green", "0")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
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
