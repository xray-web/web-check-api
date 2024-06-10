package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/checks"
	"github.com/xray-web/web-check-api/handlers/testutils"
)

func TestHandleCarbon(t *testing.T) {
	t.Parallel()

	html := `<html><body>Test</body></html>`
	client := testutils.MockClient(
		testutils.Response(http.StatusOK, []byte(html)),
		testutils.ResponseJSON(http.StatusOK, map[string]interface{}{
			"statistics": map[string]interface{}{
				"adjustedBytes": float64(len(html)),
				"energy":        0.005,
			},
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/carbon?url=http://test.com", nil)
	rec := httptest.NewRecorder()
	HandleCarbon(checks.NewCarbon(client)).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var data checks.CarbonData
	err := json.Unmarshal(rec.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.NotEmpty(t, data.ScanUrl)
	assert.Equal(t, "http://test.com", data.ScanUrl)

	assert.NotEmpty(t, data.Statistics)
	assert.Equal(t, float64(len(html)), data.Statistics.AdjustedBytes)
	assert.Equal(t, 0.005, data.Statistics.Energy)
}
