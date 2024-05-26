package controllers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DnsController struct{}

// ARecord represents an A record.
type ARecord struct {
	Address string `json:"address"`
	Family  int    `json:"family"`
}

// DNSResponse holds various DNS records.
type DNSResponse struct {
	A     ARecord  `json:"A"`
	AAAA  []string `json:"AAAA"`
	MX    []string `json:"MX"`
	TXT   []string `json:"TXT"`
	NS    []string `json:"NS"`
	CNAME []string `json:"CNAME"`
	SOA   []string `json:"SOA"`
	SRV   []string `json:"SRV"`
	PTR   []string `json:"PTR"`
}

func resolveDNSRecords(ctx context.Context, hostname string) (*DNSResponse, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second,
			}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}

	var aRecord ARecord
	var aaaaRecords, nsRecords, ptrRecords, mxRecords []string
	var txtRecords, cnameRecords, soaRecords, srvRecords []string
	var err error

	// Resolve A and AAAA records
	lookupResults, err := r.LookupIPAddr(ctx, hostname)
	if err == nil {
		for _, ip := range lookupResults {
			if ip.IP.To4() != nil {
				aRecord = ARecord{
					Address: ip.IP.String(),
					Family:  4,
				}
				aaaaRecords = append(aaaaRecords, ip.IP.String())
			} else if ip.IP.To16() != nil {
				mxRecords = append(mxRecords, ip.IP.String())
			}
		}
	}

	txtResults, _ := r.LookupTXT(ctx, hostname)
	txtRecords = append(txtRecords, txtResults...)

	nsResults, _ := r.LookupNS(ctx, hostname)
	for _, ns := range nsResults {
		nsRecords = append(nsRecords, ns.Host)
	}

	cname, _ := r.LookupCNAME(ctx, hostname)
	if cname != "" {
		cnameRecords = append(cnameRecords, cname)
	}

	// SOA records are not directly supported in Go's net package

	_, srvResults, _ := r.LookupSRV(ctx, "", "", hostname)
	for _, srv := range srvResults {
		srvRecords = append(srvRecords, srv.Target)
	}

	ptrResults, _ := r.LookupAddr(ctx, hostname)
	ptrRecords = append(ptrRecords, ptrResults...)

	return &DNSResponse{
		A:     aRecord,
		AAAA:  aaaaRecords,
		MX:    mxRecords,
		TXT:   txtRecords,
		NS:    nsRecords,
		CNAME: cnameRecords,
		SOA:   soaRecords,
		SRV:   srvRecords,
		PTR:   ptrRecords,
	}, err
}

func (ctrl *DnsController) DnsHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Extract the hostname from the URL
	hostname := url
	if strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://") {
		hostname = strings.ReplaceAll(hostname, "http://", "")
		hostname = strings.ReplaceAll(hostname, "https://", "")
		if parts := strings.Split(hostname, "/"); len(parts) > 0 {
			hostname = parts[0]
		}
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Resolve DNS records
	dnsResponse, err := resolveDNSRecords(ctx, hostname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error resolving DNS: %v", err)})
		return
	}

	// Ensure that TXT, NS, SOA, SRV, PTR are empty arrays if they are nil
	if dnsResponse.MX == nil {
		dnsResponse.MX = []string{}
	}
	if dnsResponse.TXT == nil {
		dnsResponse.TXT = []string{}
	}
	if dnsResponse.NS == nil {
		dnsResponse.NS = []string{}
	}
	if dnsResponse.SOA == nil {
		dnsResponse.SOA = []string{}
	}
	if dnsResponse.SRV == nil {
		dnsResponse.SRV = []string{}
	}
	if dnsResponse.PTR == nil {
		dnsResponse.PTR = []string{}
	}

	c.JSON(http.StatusOK, dnsResponse)
}
