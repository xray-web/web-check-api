package checks

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/testutils"
)

func TestGetLinkedPages(t *testing.T) {
	t.Parallel()
	testTargetURL := &url.URL{
		Scheme: "http",
		Host:   "internal.com",
	}
	testHTML := []byte(`
		<a href="http://internal.com/#heading"></a>
		<a href="//internal.com/1"></a>
		<a href="2"></a>
		<a href="/2"></a>
		<a href="http://external.com/1"></a>
		<a href="https://external.com/2"></a>
		<a href="http://external.com/2"></a>
		<a href="://external.com"></a>
		`)
	client := testutils.MockClient(testutils.Response(http.StatusOK, testHTML))
	actualLinkedPagesData, err := NewLinkedPages(client).GetLinkedPages(context.TODO(), testTargetURL)
	assert.NoError(t, err)
	assert.Equal(t, LinkedPagesData{
		Internal: []string{
			"http://internal.com/2",
			"http://internal.com/#heading",
			"http://internal.com/1",
		},
		External: []string{
			"http://external.com/1",
			"http://external.com/2",
			"https://external.com/2",
		},
	}, actualLinkedPagesData)
}

func TestResolveURL(t *testing.T) {
	t.Parallel()
	baseURL := url.URL{
		Scheme: "http",
		Host:   "example.com",
	}

	tests := []struct {
		name                string
		href                string
		expectedResolvedURL (*url.URL)
		expectedErrorExists bool
	}{
		{
			name:                "empty href",
			href:                "",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com"},
			expectedErrorExists: false,
		},
		{
			name:                "missing scheme",
			href:                "://example.com",
			expectedResolvedURL: nil,
			expectedErrorExists: true,
		},
		{
			name:                "relative scheme",
			href:                "//example.com",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com"},
			expectedErrorExists: false,
		},
		{
			name:                "valid absolute url without path",
			href:                "http://example.com",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com"},
			expectedErrorExists: false,
		},
		{
			name:                "valid absolute url with path",
			href:                "http://example.com/123",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/123"},
			expectedErrorExists: false,
		},
		{
			name:                "valid relative url with leading slash",
			href:                "/123",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/123"},
			expectedErrorExists: false,
		},
		{
			name:                "valid relative url without leading slash",
			href:                "123",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/123"},
			expectedErrorExists: false,
		},
		{
			name:                "valid relative url edge case",
			href:                "http//example.com",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/http//example.com"},
			expectedErrorExists: false,
		},
		{
			name:                "valid relative url fragment",
			href:                "#heading",
			expectedResolvedURL: &url.URL{Scheme: "http", Host: "example.com", Fragment: "heading"},
			expectedErrorExists: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualResolvedURL, err := resolveURL(&baseURL, tc.href)
			assert.Equal(t, tc.expectedResolvedURL, actualResolvedURL)
			if tc.expectedErrorExists {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSortURLsByFrequency(t *testing.T) {
	t.Parallel()
	testLinksMap := map[string]int{
		"https://example.com":  1,
		"https://example2.com": 2,
		"https://example3.com": 3,
	}

	expectedSortedLinks := []string{
		"https://example3.com",
		"https://example2.com",
		"https://example.com",
	}

	actualSortedLinks := sortURLsByFrequency(testLinksMap)
	assert.Equal(t, expectedSortedLinks, actualSortedLinks)
}
