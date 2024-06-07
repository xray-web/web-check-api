package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetLinks(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/get-links?url=www.google.com", nil)
	rec := httptest.NewRecorder()
	HandleGetLinks().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response LinkResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response)
	assert.NotNil(t, response.Internal)
	assert.NotNil(t, response.External)
}
