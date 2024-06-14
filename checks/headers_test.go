package checks

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/testutils"
)

func TestList(t *testing.T) {
	t.Parallel()

	c := testutils.MockClient(&http.Response{
		Header: http.Header{
			"Cache-Control":    {"private, max-age=0"},
			"X-Xss-Protection": {"0"},
		},
	})
	h := NewHeaders(c)

	actual, err := h.List(context.Background(), "example.com")
	assert.NoError(t, err)

	assert.Equal(t, "private, max-age=0", actual["Cache-Control"])
	assert.Equal(t, "0", actual["X-Xss-Protection"])
}
