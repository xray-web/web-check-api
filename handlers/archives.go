package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const archiveAPIURL = "https://web.archive.org/cdx/search/cdx"

func convertTimestampToDate(timestamp string) (time.Time, error) {
	year, err := strconv.Atoi(timestamp[0:4])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(timestamp[4:6])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(timestamp[6:8])
	if err != nil {
		return time.Time{}, err
	}
	hour, err := strconv.Atoi(timestamp[8:10])
	if err != nil {
		return time.Time{}, err
	}
	minute, err := strconv.Atoi(timestamp[10:12])
	if err != nil {
		return time.Time{}, err
	}
	second, err := strconv.Atoi(timestamp[12:14])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
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

func getScanFrequency(firstScan, lastScan time.Time, totalScans, changeCount int) map[string]float64 {
	formatToTwoDecimal := func(num float64) float64 {
		return math.Round(num*100) / 100
	}

	dayFactor := lastScan.Sub(firstScan).Hours() / 24
	daysBetweenScans := formatToTwoDecimal(dayFactor / float64(totalScans))
	daysBetweenChanges := formatToTwoDecimal(dayFactor / float64(changeCount))
	scansPerDay := formatToTwoDecimal(float64(totalScans-1) / dayFactor)
	changesPerDay := formatToTwoDecimal(float64(changeCount) / dayFactor)

	// Handle NaN values
	if math.IsNaN(daysBetweenScans) {
		daysBetweenScans = 0
	}
	if math.IsNaN(daysBetweenChanges) {
		daysBetweenChanges = 0
	}
	if math.IsNaN(scansPerDay) {
		scansPerDay = 0
	}
	if math.IsNaN(changesPerDay) {
		changesPerDay = 0
	}

	return map[string]float64{
		"daysBetweenScans":   daysBetweenScans,
		"daysBetweenChanges": daysBetweenChanges,
		"scansPerDay":        scansPerDay,
		"changesPerDay":      changesPerDay,
	}
}

func getWaybackData(url string) (map[string]interface{}, error) {
	cdxUrl := fmt.Sprintf("%s?url=%s&output=json&fl=timestamp,statuscode,digest,length,offset", archiveAPIURL, url)

	resp, err := http.Get(cdxUrl)
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

	// Remove the header row
	data = data[1:]

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
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlParam, "http://") && !strings.HasPrefix(urlParam, "https://") {
			urlParam = "http://" + urlParam
		}

		data, err := getWaybackData(urlParam)
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
