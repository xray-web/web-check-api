package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetPorts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		url          string
		expectedCode int
		expectedBody map[string][]int
	}{
		{
			url:          "open.com",
			expectedCode: http.StatusOK,
		},
		{
			url:          "closed.com",
			expectedCode: http.StatusOK,
		},
		{
			url:          "mixed.com",
			expectedCode: http.StatusOK,
		},
		{
			url:          "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.url, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/check-ports?url="+tc.url, nil)
			rec := httptest.NewRecorder()
			HandleGetPorts().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				var response map[string][]int
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				openPorts, ok1 := response["openPorts"]
				failedPorts, ok2 := response["failedPorts"]

				assert.True(t, ok1)
				assert.True(t, ok2)
				assert.NotNil(t, openPorts)
				assert.NotNil(t, failedPorts)
			} else {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "missing URL parameter")
			}
		})
	}
}
