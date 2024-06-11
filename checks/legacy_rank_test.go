package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/checks/store/legacyrank"
)

func TestLegacyRank(t *testing.T) {
	t.Parallel()

	t.Run("get rank", func(t *testing.T) {
		t.Parallel()
		lr := NewLegacyRank(legacyrank.GetterFunc(func(domain string) (int, error) {
			return 1, nil
		}))
		dr, err := lr.LegacyRank("example.com")
		assert.NoError(t, err)
		assert.Equal(t, 1, dr.Rank)
		assert.Equal(t, "example.com", dr.Domain)
	})
}
