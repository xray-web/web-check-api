package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractURL(t *testing.T) {
	// NOTE: due to golangs strict checking of URLs, it is hard to trigger the invalid URL error
	// so we will only test the missing URL parameter and the valid URL cases and use external tools to test the invalid URL case
	t.Parallel()

	t.Run("missing URL parameter", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "http://test.com", nil)
		u, err := extractURL(req)
		assert.Equal(t, ErrMissingURLParameter, err)
		assert.Nil(t, u)
	})

	t.Run("valid URL without scheme", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "http://test.com?url=example.com", nil)
		u, err := extractURL(req)
		assert.Nil(t, err)
		assert.Equal(t, "http", u.Scheme)
		assert.Equal(t, "http://example.com", u.String())
		assert.Equal(t, "example.com", u.Host)
		assert.Equal(t, "example.com", u.Hostname())
	})

	t.Run("valid URL with scheme", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "http://test.com?url=http://example.com", nil)
		u, err := extractURL(req)
		assert.Nil(t, err)
		assert.Equal(t, "http", u.Scheme)
		assert.Equal(t, "http://example.com", u.String())
		assert.Equal(t, "example.com", u.Host)
		assert.Equal(t, "example.com", u.Hostname())
	})
}
