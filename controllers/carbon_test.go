package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCarbonHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	ctrl := &controllers.CarbonController{}

	r.GET("/carbon", ctrl.CarbonHandler)

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

	req, _ := http.NewRequest(http.MethodGet, "/carbon?url=http://test.com", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, but got %d", w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error unmarshaling response body")

	assert.NotNil(t, response["scanUrl"], "scanUrl should not be nil")
	assert.Equal(t, "http://test.com", response["scanUrl"], "Expected scanUrl to be 'http://test.com', but got %v", response["scanUrl"])

	assert.NotNil(t, response["statistics"], "statistics should not be nil")
	stats, ok := response["statistics"].(map[string]interface{})
	assert.True(t, ok, "Expected statistics to be a map, but got %T", response["statistics"])
	assert.Equal(t, float64(htmlSize), stats["adjustedBytes"], "Expected adjustedBytes to be %d, but got %v", htmlSize, stats["adjustedBytes"])
	assert.Equal(t, 0.005, stats["energy"], "Expected energy to be 0.005, but got %v", stats["energy"])
}
