package checks

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/testutils"
)

func TestCarbonHtmlSize(t *testing.T) {
	t.Parallel()
	const htmlBody = `<html><body>Test</body></html>`
	var size = len(htmlBody)
	client := testutils.MockClient(testutils.Response(http.StatusOK, []byte(htmlBody)))
	c := NewCarbon(client)
	size, err := c.HtmlSize(context.TODO(), "/carbon")
	assert.NoError(t, err)
	assert.Equal(t, len(htmlBody), size)
}
