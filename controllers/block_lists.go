package controllers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type BlockListsController struct{}

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

func isDomainBlocked(domain, serverIP string) bool {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 3,
			}
			return d.DialContext(ctx, network, serverIP+":53")
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ips, err := resolver.LookupIP(ctx, "ip4", domain)
	if err != nil {
		// if there's an error, consider it not blocked
		return false
	}

	return slices.ContainsFunc(ips, func(ip net.IP) bool {
		return slices.Contains(knownBlockIPs, ip.String())
	})
}

func checkDomainAgainstDNSServers(domain string) []Blocklist {
	var lock sync.Mutex
	var wg sync.WaitGroup
	limit := make(chan struct{}, 5)

	var results []Blocklist

	for _, server := range DNS_SERVERS {
		wg.Add(1)
		go func(server dnsServer) {
			limit <- struct{}{}
			defer func() {
				<-limit
				wg.Done()
			}()

			isBlocked := isDomainBlocked(domain, server.IP)
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

	return results
}

func (ctrl *BlockListsController) BlockListsHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL parameter"})
		return
	}

	domain, err := urlToDomain(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	results := checkDomainAgainstDNSServers(domain)
	c.JSON(http.StatusOK, gin.H{"blocklists": results})
}

func urlToDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}

func HandleBlockLists() http.Handler {
	type Response struct {
		BlockLists []Blocklist `json:"blocklists"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")
		if rawURL == "" {
			JSONError(w, "Missing URL parameter", http.StatusBadRequest)
			return
		}
		domain, err := urlToDomain(rawURL)
		if err != nil {
			JSONError(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(Response{BlockLists: checkDomainAgainstDNSServers(domain)})
	})
}
