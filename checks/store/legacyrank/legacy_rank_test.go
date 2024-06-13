package legacyrank_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/checks/store/legacyrank"
)

func TestInMemoryStore(t *testing.T) {
	t.Parallel()

	t.Run("get google rank", func(t *testing.T) {
		t.Parallel()
		ims := legacyrank.NewInMemoryStore()
		dr, err := ims.GetLegacyRank("google.com")
		assert.NoError(t, err, dr)
	})

	t.Run("get microsoft rank", func(t *testing.T) {
		t.Parallel()
		ims := legacyrank.NewInMemoryStore()
		dr, err := ims.GetLegacyRank("microsoft.com")
		assert.NoError(t, err, dr)
	})
}
