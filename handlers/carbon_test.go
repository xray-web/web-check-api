package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestHandleCarbon(t *testing.T) {
	// t.Parallel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test.com",
		httpmock.NewStringResponder(200, "<html><body>Test</body></html>"))

	htmlSize := len("<html><body>Test</body></html>")

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.websitecarbon.com/data?bytes=%d&green=0", htmlSize),
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"statistics": map[string]interface{}{
				"adjustedBytes": float64(htmlSize),
				"energy":        0.005,
			},
		}))

	req := httptest.NewRequest(http.MethodGet, "/carbon?url=http://test.com", nil)
	rec := httptest.NewRecorder()
	HandleCarbon().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200, but got %d", rec.Code)

	var data CarbonData
	err := json.Unmarshal(rec.Body.Bytes(), &data)
	assert.NoError(t, err, "Error unmarshaling response body")

	assert.NotEmpty(t, data.ScanUrl, "scanUrl should not be nil")
	assert.Equal(t, "http://test.com", data.ScanUrl, "Expected scanUrl to be 'http://test.com', but got %v", data.ScanUrl)

	assert.NotEmpty(t, data.Statistics, "statistics should not be nil")
	stats := data.Statistics
	assert.Equal(t, float64(htmlSize), stats.AdjustedBytes, "Expected adjustedBytes to be %d, but got %v", htmlSize, stats.AdjustedBytes)
	assert.Equal(t, 0.005, stats.Energy, "Expected energy to be 0.005, but got %v", stats.Energy)
}
