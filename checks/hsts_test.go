package checks

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/testutils"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	t.Run("given an empty header", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(&http.Response{
			Header: http.Header{StrictTransportSecurity: []string{""}}})
		h := NewHsts(client)

		actual, err := h.Validate(context.Background(), "test.com")
		assert.NoError(t, err)

		assert.Equal(t, NilHeadersError, actual.Message)
		assert.False(t, actual.Compatible)
		assert.Empty(t, actual.HSTSHeader)
	})

	t.Run("given a header without max age", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(&http.Response{
			Header: http.Header{StrictTransportSecurity: []string{"includeSubDomains; preload"}}})
		h := NewHsts(client)

		actual, err := h.Validate(context.Background(), "test.com")
		assert.NoError(t, err)

		assert.Equal(t, MaxAgeError, actual.Message)
		assert.False(t, actual.Compatible)
		assert.Empty(t, actual.HSTSHeader)
	})

	t.Run("given max age less than 10886400", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(&http.Response{
			Header: http.Header{StrictTransportSecurity: []string{"max-age=47; includeSubDomains; preload"}}})
		h := NewHsts(client)

		actual, err := h.Validate(context.Background(), "test.com")
		assert.NoError(t, err)

		assert.Equal(t, MaxAgeError, actual.Message)
		assert.False(t, actual.Compatible)
		assert.Empty(t, actual.HSTSHeader)
	})

	t.Run("given a header without includeSubDomains", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(&http.Response{
			Header: http.Header{StrictTransportSecurity: []string{"max-age=47474747; preload"}}})
		h := NewHsts(client)

		actual, err := h.Validate(context.Background(), "test.com")
		assert.NoError(t, err)

		assert.Equal(t, SubdomainsError, actual.Message)
		assert.False(t, actual.Compatible)
		assert.Empty(t, actual.HSTSHeader)
	})

	t.Run("given a header without preload", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(&http.Response{
			Header: http.Header{StrictTransportSecurity: []string{"max-age=47474747; includeSubDomains"}}})
		h := NewHsts(client)

		actual, err := h.Validate(context.Background(), "test.com")
		assert.NoError(t, err)

		assert.Equal(t, PreloadError, actual.Message)
		assert.False(t, actual.Compatible)
		assert.Empty(t, actual.HSTSHeader)
	})

	t.Run("given a valid header", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(&http.Response{
			Header: http.Header{StrictTransportSecurity: []string{"max-age=47474747; includeSubDomains; preload"}}})
		h := NewHsts(client)

		actual, err := h.Validate(context.Background(), "test.com")
		assert.NoError(t, err)

		assert.Equal(t, HstsSuccess, actual.Message)
		assert.True(t, actual.Compatible)
		assert.NotEmpty(t, actual.HSTSHeader)
	})
}

func TestExtractMaxAgeFromHeader(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		header   string
		expected string
	}{
		{"give valid header", "max-age=47474747; includeSubDomains; preload", "47474747"},
		{"given an empty header", "", ""},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := extractMaxAgeFromHeader(tc.header)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
