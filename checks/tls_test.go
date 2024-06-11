package checks

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/testutils"
)

func TestTLS(t *testing.T) {
	t.Parallel()

	t.Run("Valid URL with successful scan", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(
			testutils.Response(http.StatusOK, []byte(`{"scan_id": 12345}`)),
			testutils.Response(http.StatusOK, []byte(`{"grade": "A+"}`)),
		)

		tls, err := NewTls(client).GetScanResults(context.TODO(), "example.com")
		assert.NoError(t, err)
		assert.Equal(t, "A+", tls["grade"])
	})
}
