package checks

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/testutils"
)

func TestSocialTagsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		s := SocialTagsData{}
		assert.True(t, s.Empty())
	})

	t.Run("Not empty", func(t *testing.T) {
		t.Parallel()

		s := SocialTagsData{
			Title: "Example Domain",
		}
		assert.False(t, s.Empty())
	})
}

func TestNewSocialTags(t *testing.T) {
	t.Parallel()

	t.Run("No social tags", func(t *testing.T) {
		t.Parallel()

		client := testutils.MockClient(testutils.Response(http.StatusOK, []byte{}))
		tags, err := NewSocialTags(client).GetSocialTags(context.TODO(), "http://example.com")
		assert.NoError(t, err)
		assert.True(t, tags.Empty())
	})

	t.Run("Social tags", func(t *testing.T) {
		t.Parallel()

		var html = []byte(`
		<html>
			<head>
				<title>Example Domain</title>
				<meta name="description" content="Example description">
				<meta property="og:title" content="Example OG Title">
			</head>
			<body></body>
		</html>
		`)
		client := testutils.MockClient(testutils.Response(http.StatusOK, html))
		tags, err := NewSocialTags(client).GetSocialTags(context.TODO(), "http://example.com")
		assert.NoError(t, err)
		assert.False(t, tags.Empty())
		assert.Equal(t, "Example description", tags.Description)
		assert.Equal(t, "Example Domain", tags.Title)
		assert.Equal(t, "Example OG Title", tags.OgTitle)
	})
}
