package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const archiveAPIURL = "https://web.archive.org/cdx/search/cdx"

func convertTimestampToDate(timestamp string) (time.Time, error) {
	mask := "20060102150405"
	return time.Parse(mask, timestamp)
}

func countPageChanges(results [][]string) int {
	prevDigest := ""
	changeCount := -1
	for _, curr := range results {
		if curr[2] != prevDigest {
			prevDigest = curr[2]
			changeCount++
		}
	}
	return changeCount
}

func getAveragePageSize(scans [][]string) int {
	totalSize := 0
	for _, scan := range scans {
		size, err := strconv.Atoi(scan[3])
		if err != nil {
			continue
		}
		totalSize += size
	}
	return totalSize / len(scans)
}

func getScanFrequency(firstScan, lastScan time.Time, totalScans, changeCount int) map[string]string {
	formatToTwoDecimal := func(num float64) string {
		return fmt.Sprintf("%.2f", num)
	}

	dayFactor := lastScan.Sub(firstScan).Hours() / 24
	daysBetweenScans := formatToTwoDecimal(dayFactor / float64(totalScans))
	daysBetweenChanges := formatToTwoDecimal(dayFactor / float64(changeCount))
	scansPerDay := formatToTwoDecimal(float64(totalScans-1) / dayFactor)
	changesPerDay := formatToTwoDecimal(float64(changeCount) / dayFactor)

	if math.IsNaN(dayFactor / float64(totalScans)) {
		daysBetweenScans = "0.00"
	}
	if math.IsNaN(dayFactor / float64(changeCount)) {
		daysBetweenChanges = "0.00"
	}
	if math.IsNaN(float64(totalScans-1) / dayFactor) {
		scansPerDay = "0.00"
	}
	if math.IsNaN(float64(changeCount) / dayFactor) {
		changesPerDay = "0.00"
	}

	return map[string]string{
		"daysBetweenScans":   daysBetweenScans,
		"daysBetweenChanges": daysBetweenChanges,
		"scansPerDay":        scansPerDay,
		"changesPerDay":      changesPerDay,
	}
}

func getWaybackData(url *url.URL) (map[string]interface{}, error) {
	cdxUrl := fmt.Sprintf("%s?url=%s&output=json&fl=timestamp,statuscode,digest,length,offset", archiveAPIURL, url)

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(cdxUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data [][]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data) <= 1 {
		return map[string]interface{}{
			"skipped": "Site has never before been archived via the Wayback Machine",
		}, nil
	}

	if len(data) < 1 {
		return nil, fmt.Errorf("data slice is empty")
	}

	// Remove the header row
	data = data[1:]

	if len(data) < 1 {
		return nil, fmt.Errorf("data slice became empty after removing the first element")
	}

	// Access the first element of the remaining data
	firstScan, err := convertTimestampToDate(data[0][0])
	if err != nil {
		return nil, err
	}
	lastScan, err := convertTimestampToDate(data[len(data)-1][0])
	if err != nil {
		return nil, err
	}
	totalScans := len(data)
	changeCount := countPageChanges(data)

	return map[string]interface{}{
		"firstScan":       firstScan.Format(time.RFC3339),
		"lastScan":        lastScan.Format(time.RFC3339),
		"totalScans":      totalScans,
		"changeCount":     changeCount,
		"averagePageSize": getAveragePageSize(data),
		"scanFrequency":   getScanFrequency(firstScan, lastScan, totalScans, changeCount),
		"scans":           data,
		"scanUrl":         url,
	}, nil
}

func HandleArchives() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		data, err := getWaybackData(rawURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching Wayback data: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	})
}
