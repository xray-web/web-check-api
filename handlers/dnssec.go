package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const dnsGoogleURL = "https://dns.google/resolve"

func resolveDNS(domain string) (map[string]interface{}, error) {
	dnsTypes := []string{"DNSKEY", "DS", "RRSIG"}
	records := make(map[string]interface{})
	domain = strings.TrimSuffix(domain, "")
	for _, typ := range dnsTypes {

		url := fmt.Sprintf("%s?name=%s&type=%s", dnsGoogleURL, url.PathEscape(domain), typ)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching %s record: %s", typ, err.Error())
		}
		defer resp.Body.Close()

		var dnsResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&dnsResponse); err != nil {
			return nil, fmt.Errorf("error decoding JSON response for %s record: %s", typ, err.Error())
		}

		// Extract comments from the DNS response
		comment := ""
		if comments, ok := dnsResponse["Comment"]; ok {
			comment = comments.(string)
		}

		if answer, ok := dnsResponse["Answer"]; ok {
			records[typ] = map[string]interface{}{
				"isFound":  true,
				"answer":   answer,
				"response": dnsResponse,
				"Comment":  comment, // Include comment in the output
			}
		} else {
			records[typ] = map[string]interface{}{
				"isFound":  false,
				"answer":   nil,
				"response": dnsResponse,
				"Comment":  comment, // Include comment in the output
			}
		}
	}

	return records, nil
}

func HandleDnsSec() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("url")
		if domain == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		records, err := resolveDNS(domain)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, records, http.StatusOK)
	})
}
