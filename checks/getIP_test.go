package checks

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	ipAddresses := []IpAddress{
		{net.ParseIP("216.58.201.110"), 4},
		{net.ParseIP("2a00:1450:4009:826::200e"), 6},
	}
	i := NewIp(IpGetterFunc(func(ctx context.Context, host string) ([]IpAddress, error) {
		return ipAddresses, nil
	}))
	actual, err := i.Lookup(context.Background(), "google.com")
	assert.NoError(t, err)

	assert.Equal(t, ipAddresses[0].Address, actual[0].Address)
	assert.Equal(t, 4, actual[0].Family)

	assert.Equal(t, ipAddresses[1].Address, actual[1].Address)
	assert.Equal(t, 6, actual[1].Family)
}
