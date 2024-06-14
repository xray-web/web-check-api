package checks

import (
	"context"
	"net"
)

type IpAddress struct {
	Address net.IP `json:"ip"`
	Family  int    `json:"family"`
}

type IpGetter interface {
	GetIp(ctx context.Context, host string) ([]IpAddress, error)
}

type IpGetterFunc func(ctx context.Context, host string) ([]IpAddress, error)

func (f IpGetterFunc) GetIp(ctx context.Context, host string) ([]IpAddress, error) {
	return f(ctx, host)
}

type NetIp struct{}

func NewNetIp() *NetIp {
	return &NetIp{}
}

func (l *NetIp) GetIp(ctx context.Context, host string) ([]IpAddress, error) {
	resolver := &net.Resolver{
		PreferGo: true,
	}
	ip4, err := resolver.LookupIP(ctx, "ip4", host)
	if err != nil {
		return nil, err
	}
	ip6, err := resolver.LookupIP(ctx, "ip6", host)
	if err != nil {
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

type Ip struct {
	getter IpGetter
}

func NewIp(l IpGetter) *Ip {
	return &Ip{getter: l}
}

func (i *Ip) Lookup(ctx context.Context, host string) ([]IpAddress, error) {
	return i.getter.GetIp(ctx, host)
}
