package checks

import (
	"context"
	"net"

	"github.com/xray-web/web-check-api/checks/clients/ip"
)

type IpAddress struct {
	Address net.IP `json:"ip"`
	Family  int    `json:"family"`
}

type NetIp struct {
	lookup ip.Lookup
}

func NewNetIp(lookup ip.Lookup) *NetIp {
	return &NetIp{lookup: lookup}
}

func (l *NetIp) GetIp(ctx context.Context, host string) ([]IpAddress, error) {
	ip4, err := l.lookup.LookupIP(ctx, "ip4", host)
	if err != nil {
		// do nothing
	}
	ip6, err := l.lookup.LookupIP(ctx, "ip6", host)
	if err != nil && len(ip4) == 0 && len(ip6) == 0 {
		return nil, err
	}

	var ipAddresses []IpAddress
	for _, ip := range ip4 {
		ipAddresses = append(ipAddresses, IpAddress{Address: ip, Family: 4})
	}
	for _, ip := range ip6 {
		ipAddresses = append(ipAddresses, IpAddress{Address: ip, Family: 6})
	}

	return ipAddresses, nil
}
