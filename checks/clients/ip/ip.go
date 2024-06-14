package ip

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Lookup interface {
	LookupIP(ctx context.Context, network string, host string) ([]net.IP, error)
}

type LookupFunc func(ctx context.Context, network string, host string) ([]net.IP, error)

func (fn LookupFunc) LookupIP(ctx context.Context, network string, host string) ([]net.IP, error) {
	return fn(ctx, network, host)
}

// NetLookup is a client for looking up IP addresses using a net.Resolver.
type NetLookup struct{}

func (l *NetLookup) LookupIP(ctx context.Context, network string, host string) ([]net.IP, error) {
	netResolver := &net.Resolver{
		PreferGo: true,
	}
	return netResolver.LookupIP(ctx, network, host)
}

type DNSLookup interface {
	DNSLookupIP(ctx context.Context, network, host, dns string) ([]net.IP, error)
}

type DNSLookupFunc func(ctx context.Context, network, host, dns string) ([]net.IP, error)

func (fn DNSLookupFunc) DNSLookupIP(ctx context.Context, network, host, dns string) ([]net.IP, error) {
	return fn(ctx, network, host, dns)
}

// DNSLookup is a client for looking up IP addresses with a custom DNS server.
type NetDNSLookup struct{}

func (l *NetDNSLookup) DNSLookupIP(ctx context.Context, network, host, dns string) ([]net.IP, error) {
	netResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 3 * time.Second,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", dns, 53))
		},
	}
	return netResolver.LookupIP(ctx, network, host)
}
