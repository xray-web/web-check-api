package controllers

import (
	"net"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type BlockListsController struct{}

var DNS_SERVERS = []struct {
	Name string
	IP   string
}{
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
	ips, err := net.LookupIP(domain)
	if err != nil {
		return true
	}

	for _, ip := range ips {
		if contains(knownBlockIPs, ip.String()) {
			return true
		}
	}

	return false
}

func checkDomainAgainstDNSServers(domain string) []Blocklist {
	var results []Blocklist

	for _, server := range DNS_SERVERS {
		isBlocked := isDomainBlocked(domain, server.IP)
		results = append(results, Blocklist{
			Server:    server.Name,
			ServerIP:  server.IP,
			IsBlocked: isBlocked,
		})
	}

	return results
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func (ctrl *BlockListsController) BlockListsHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL parameter"})
		return
	}

	domain, err := urlToDomain(url)
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
