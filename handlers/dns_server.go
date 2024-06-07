package handlers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Result holds the information for each resolved address.
type Result struct {
	Address           string   `json:"address"`
	Hostname          []string `json:"hostname"`
	DOHDirectSupports bool     `json:"dohDirectSupports"`
}

// Response holds the final response structure.
type Response struct {
	Domain string   `json:"domain"`
	DNS    []Result `json:"dns"`
}

func resolveDNSServer(ctx context.Context, domain string) ([]Result, error) {
	addrs, err := net.DefaultResolver.LookupIP(ctx, "ip4", domain)
	if err != nil {
		return nil, fmt.Errorf("could not resolve DNS: %v", err)
	}

	var results []Result
	for _, addr := range addrs {
		ip := addr.String()

		hostnames, err := net.DefaultResolver.LookupAddr(ctx, ip)
		if err != nil {
			hostnames = nil
		}

		dohDirectSupports := false
		client := http.Client{Timeout: 3 * time.Second}
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s/dns-query", ip), nil)
		if err == nil {
			resp, err := client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				dohDirectSupports = true
				resp.Body.Close()
			}
		}

		result := Result{
			Address:           ip,
			Hostname:          hostnames,
			DOHDirectSupports: dohDirectSupports,
		}
		results = append(results, result)
	}
	return results, nil
}

func HandleDNSServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		// Extract the domain from the URL
		domain := strings.ReplaceAll(url, "http://", "")
		domain = strings.ReplaceAll(domain, "https://", "")
		domain = strings.TrimSuffix(domain, "/")

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		results, err := resolveDNSServer(ctx, domain)
		if err != nil {
			JSONError(w, fmt.Errorf("error resolving DNS: %v", err), http.StatusInternalServerError)
			return
		}

		response := Response{
			Domain: domain,
			DNS:    results,
		}

		JSON(w, response, http.StatusOK)
	})
}
