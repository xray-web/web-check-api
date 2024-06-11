package checks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Tls struct {
	client *http.Client
}

func NewTls(client *http.Client) *Tls {
	return &Tls{client: client}
}

func (t *Tls) initiateScan(ctx context.Context, domain string) (int, error) {
	const scanUrl = "https://tls-observatory.services.mozilla.com/api/v1/scan"

	formData := url.Values{"target": {domain}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, scanUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var res struct {
		ScanID int `json:"scan_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return -1, err
	}

	if res.ScanID == 0 {
		return res.ScanID, errors.New("failed to get scan_id from TLS Observatory")
	}

	return res.ScanID, nil
}

func (t *Tls) GetScanResults(ctx context.Context, domain string) (map[string]interface{}, error) {
	scanID, err := t.initiateScan(ctx, domain)
	if err != nil {
		return nil, err
	}

	const scanUrl = "https://tls-observatory.services.mozilla.com/api/v1/results"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, scanUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id", strconv.Itoa(scanID))
	req.URL.RawQuery = q.Encode()

	resp, err := t.client.Do(req)
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
