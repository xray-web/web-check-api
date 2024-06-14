package checks

import (
	"context"
	"net"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/xray-web/web-check-api/checks/clients/ip"
)

type dnsServer struct {
	Name string
	IP   string
}

var DNS_SERVERS = []dnsServer{
	{Name: "AdGuard", IP: "176.103.130.130"},
	{Name: "AdGuard Family", IP: "176.103.130.132"},
	{Name: "CleanBrowsing Adult", IP: "185.228.168.10"},
	{Name: "CleanBrowsing Family", IP: "185.228.168.168"},
	{Name: "CleanBrowsing Security", IP: "185.228.168.9"},
	{Name: "CloudFlare", IP: "1.1.1.1"},
	{Name: "CloudFlare Family", IP: "1.1.1.3"},
	{Name: "Comodo Secure", IP: "8.26.56.26"},
	{Name: "Google DNS", IP: "8.8.8.8"},
	{Name: "Neustar Family", IP: "156.154.70.3"},
	{Name: "Neustar Protection", IP: "156.154.70.2"},
	{Name: "Norton Family", IP: "199.85.126.20"},
	{Name: "OpenDNS", IP: "208.67.222.222"},
	{Name: "OpenDNS Family", IP: "208.67.222.123"},
	{Name: "Quad9", IP: "9.9.9.9"},
	{Name: "Yandex Family", IP: "77.88.8.7"},
	{Name: "Yandex Safe", IP: "77.88.8.88"},
}

var knownBlockIPs = []string{
	"146.112.61.106",
	"185.228.168.10",
	"8.26.56.26",
	"9.9.9.9",
	"208.69.38.170",
	"208.69.39.170",
	"208.67.222.222",
	"208.67.222.123",
	"199.85.126.10",
	"199.85.126.20",
	"156.154.70.22",
	"77.88.8.7",
	"77.88.8.8",
	"::1",
	"2a02:6b8::feed:0ff",
	"2a02:6b8::feed:bad",
	"2a02:6b8::feed:a11",
	"2620:119:35::35",
	"2620:119:53::53",
	"2606:4700:4700::1111",
	"2606:4700:4700::1001",
	"2001:4860:4860::8888",
	"2a0d:2a00:1::",
	"2a0d:2a00:2::",
}

type Blocklist struct {
	Server    string `json:"server"`
	ServerIP  string `json:"serverIp"`
	IsBlocked bool   `json:"isBlocked"`
}

type BlockList struct {
	lookup ip.DNSLookup
}

func NewBlockList(lookup ip.DNSLookup) *BlockList {
	return &BlockList{lookup: lookup}
}

func (b *BlockList) domainBlocked(ctx context.Context, domain, serverIP string) bool {
	ips, err := b.lookup.DNSLookupIP(ctx, "ip4", domain, serverIP)
	if err != nil {
		// if there's an error, consider it not blocked
		// TODO: return more detailed errors for each server
		return false
	}

	return slices.ContainsFunc(ips, func(ip net.IP) bool {
		return slices.Contains(knownBlockIPs, ip.String())
	})
}

func (b *BlockList) BlockedServers(ctx context.Context, domain string) []Blocklist {
	var lock sync.Mutex
	var wg sync.WaitGroup
	limit := make(chan struct{}, 5)

	var results []Blocklist

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for _, server := range DNS_SERVERS {
		wg.Add(1)
		go func(server dnsServer) {
			limit <- struct{}{}
			defer func() {
				<-limit
				wg.Done()
			}()

			isBlocked := b.domainBlocked(ctx, domain, server.IP)
			lock.Lock()
			defer lock.Unlock()
			results = append(results, Blocklist{
				Server:    server.Name,
				ServerIP:  server.IP,
				IsBlocked: isBlocked,
			})
		}(server)
	}
	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Server < results[j].Server
	})
	return results
}
