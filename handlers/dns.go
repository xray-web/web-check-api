package handlers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// ARecord represents an A record.
type ARecord struct {
	Address string `json:"address"`
	Family  int    `json:"family"`
}

// DNSResponse holds various DNS records.
type DNSResponse struct {
	A     []ARecord `json:"A"`
	AAAA  []string  `json:"AAAA"`
	MX    []string  `json:"MX"`
	TXT   []string  `json:"TXT"`
	NS    []string  `json:"NS"`
	CNAME []string  `json:"CNAME"`
	SOA   string    `json:"SOA"`
	SRV   []string  `json:"SRV"`
	PTR   []string  `json:"PTR"`
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

	var aRecords []ARecord
	var aaaaRecords, nsRecords, ptrRecords, mxRecords []string
	var txtRecords, cnameRecords, srvRecords []string
	var soaRecord string
	var err error

	// Resolve A and AAAA records
	lookupResults, err := r.LookupIPAddr(ctx, hostname)
	if err == nil {
		for _, ip := range lookupResults {
			if ip.IP.To4() != nil {
				aRecords = append(aRecords, ARecord{
					Address: ip.IP.String(),
					Family:  4,
				})
			} else if ip.IP.To16() != nil {
				aaaaRecords = append(aaaaRecords, ip.IP.String())
			}
		}
	}

	// Resolve MX records
	mxResults, _ := r.LookupMX(ctx, hostname)
	for _, mx := range mxResults {
		mxRecords = append(mxRecords, fmt.Sprintf("%s %d", mx.Host, mx.Pref))
	}

	// Resolve TXT records
	txtResults, _ := r.LookupTXT(ctx, hostname)
	txtRecords = append(txtRecords, txtResults...)

	// Resolve NS records
	nsResults, _ := r.LookupNS(ctx, hostname)
	for _, ns := range nsResults {
		nsRecords = append(nsRecords, ns.Host)
	}

	// Resolve CNAME record
	cname, _ := r.LookupCNAME(ctx, hostname)
	if cname != "" {
		cnameRecords = append(cnameRecords, cname)
	}

	// Resolve SRV records
	_, srvResults, _ := r.LookupSRV(ctx, "", "", hostname)
	for _, srv := range srvResults {
		srvRecords = append(srvRecords, fmt.Sprintf("%s %d %d %d", srv.Target, srv.Port, srv.Priority, srv.Weight))
	}

	// Resolve PTR records
	ptrResults, _ := r.LookupAddr(ctx, hostname)
	ptrRecords = append(ptrRecords, ptrResults...)

	return &DNSResponse{
		A:     aRecords,
		AAAA:  aaaaRecords,
		MX:    mxRecords,
		TXT:   txtRecords,
		NS:    nsRecords,
		CNAME: cnameRecords,
		SOA:   soaRecord,
		SRV:   srvRecords,
		PTR:   ptrRecords,
	}, err
}

func HandleDNS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Resolve DNS records
		dnsResponse, err := resolveDNSRecords(ctx, rawURL.Hostname())
		if err != nil {
			JSONError(w, fmt.Errorf("error resolving DNS: %v", err), http.StatusInternalServerError)
			return
		}

		JSON(w, dnsResponse, http.StatusOK)
	})
}
