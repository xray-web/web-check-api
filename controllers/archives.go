package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArchivesController struct{}

// Converts a timestamp string to a Go time.Time object
func convertTimestampToDate(timestamp string) time.Time {
	year, _ := strconv.Atoi(timestamp[0:4])
	month, _ := strconv.Atoi(timestamp[4:6])
	day, _ := strconv.Atoi(timestamp[6:8])
	hour, _ := strconv.Atoi(timestamp[8:10])
	minute, _ := strconv.Atoi(timestamp[10:12])
	second, _ := strconv.Atoi(timestamp[12:14])

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}

// Counts the number of changes in the page based on the digest field
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

// Calculates the average page size
func getAveragePageSize(scans [][]string) int {
	totalSize := 0

	for _, scan := range scans {
		size, _ := strconv.Atoi(scan[3])
		totalSize += size
	}

	return totalSize / len(scans)
}

// Computes various frequency metrics
func getScanFrequency(firstScan, lastScan time.Time, totalScans, changeCount int) map[string]float64 {
	dayFactor := lastScan.Sub(firstScan).Hours() / 24
	daysBetweenScans := dayFactor / float64(totalScans)
	daysBetweenChanges := dayFactor / float64(changeCount)
	scansPerDay := float64(totalScans-1) / dayFactor
	changesPerDay := float64(changeCount) / dayFactor

	return map[string]float64{
		"daysBetweenScans":   math.Round(daysBetweenScans*100) / 100,
		"daysBetweenChanges": math.Round(daysBetweenChanges*100) / 100,
		"scansPerDay":        math.Round(scansPerDay*100) / 100,
		"changesPerDay":      math.Round(changesPerDay*100) / 100,
	}
}

// Fetches data from the Wayback Machine API and processes it
func getWaybackData(url string) (map[string]interface{}, error) {
	cdxUrl := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s&output=json&fl=timestamp,statuscode,digest,length,offset", url)

	resp, err := http.Get(cdxUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching Wayback data: %v", err)
	}
	defer resp.Body.Close()

	var data [][]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	if len(data) <= 1 {
		return map[string]interface{}{"skipped": "Site has never before been archived via the Wayback Machine"}, nil
	}

	// Remove the header row
	data = data[1:]

	firstScan := convertTimestampToDate(data[0][0])
	lastScan := convertTimestampToDate(data[len(data)-1][0])
	totalScans := len(data)
	changeCount := countPageChanges(data)

	return map[string]interface{}{
		"firstScan":       firstScan,
		"lastScan":        lastScan,
		"totalScans":      totalScans,
		"changeCount":     changeCount,
		"averagePageSize": getAveragePageSize(data),
		"scanFrequency":   getScanFrequency(firstScan, lastScan, totalScans, changeCount),
		"scans":           data,
		"scanUrl":         url,
	}, nil
}

// Handler for the /archives endpoint
func (ctrl *ArchivesController) ArchivesHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		return
	}

	result, err := getWaybackData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, result)
}
