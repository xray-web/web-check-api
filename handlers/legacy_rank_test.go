package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleLegacyRank(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/legacy-rank?url=www.google.com", nil)
	rec := httptest.NewRecorder()
	HandleLegacyRank().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response RankResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.Equal(t, "www.google.com", response.Domain)

	assert.True(t, response.IsFound || !response.IsFound)
}
