package checks

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/checks/clients/ip"
)

func TestBlockList(t *testing.T) {
	t.Parallel()

	t.Run("blocked IP", func(t *testing.T) {
		t.Parallel()

		dnsLookup := ip.DNSLookupFunc(func(ctx context.Context, network, host, dns string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("146.112.61.106")}, nil
		})
		list := NewBlockList(dnsLookup).BlockedServers(context.Background(), "example.com")
		assert.Contains(t, list, Blocklist{Server: "AdGuard", ServerIP: "176.103.130.130", IsBlocked: true})
	})
}
