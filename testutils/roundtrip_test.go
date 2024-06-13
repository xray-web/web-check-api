package testutils

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTrip(t *testing.T) {
	t.Parallel()
	client := MockClient(
		Response(http.StatusOK, []byte("Hello, World!")),
		Response(http.StatusInternalServerError, []byte("Goodbye, World!")),
	)

	resp, err := client.Get("http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(b))

	resp, err = client.Get("http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	b, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Goodbye, World!", string(b))
}
