package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

const internicHostname = "whois.internic.net"
const myAPIURL = "https://whois-api-zeta.vercel.app/"

func getBaseDomain(url string) (string, error) {
	protocol := ""
	if strings.HasPrefix(url, "http://") {
		protocol = "http://"
	} else if strings.HasPrefix(url, "https://") {
		protocol = "https://"
	}
	noProtocolURL := strings.Replace(url, protocol, "", 1)
	parsed, err := parseDomain(noProtocolURL)
	if err != nil {
		return "", err
	}
	return protocol + parsed, nil
}
func parseDomain(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		domainParts := strings.Split(parts[0], ".")
		if len(domainParts) < 2 {
			return "", errors.New("invalid URL")
		}
		return domainParts[len(domainParts)-2] + "." + domainParts[len(domainParts)-1], nil
	}
	return "", errors.New("invalid URL")
}

func parseWhoisData(data string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(data, "\r\n")

	var lastKey string

	for _, line := range lines {
		index := strings.Index(line, ":")
		if index == -1 {
			if lastKey != "" {
				result[lastKey] += " " + strings.TrimSpace(line)
			}
			continue
		}
		key := strings.TrimSpace(line[:index])
		value := strings.TrimSpace(line[index+1:])
		if len(value) == 0 {
			continue
		}
		key = regexp.MustCompile(`\W+`).ReplaceAllString(key, "_")
		lastKey = key

		result[key] = value
	}

	return result
}

func fetchFromInternic(hostname string) (map[string]string, error) {
	conn, err := net.Dial("tcp", internicHostname+":43")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(hostname + "\r\n"))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(conn)
	if err != nil {
		return nil, err
	}

	parsedData := parseWhoisData(buf.String())

	if _, ok := parsedData["No_match_for"]; ok {
		return nil, errors.New("No matches found for domain in internic database")
	}

	return parsedData, nil
}

func fetchFromMyAPI(hostname string) (map[string]interface{}, error) {
	resp, err := http.Post(myAPIURL, "application/json", strings.NewReader(fmt.Sprintf(`{"domain": "%s"}`, hostname)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Log the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func HandleWhois() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
			return
		}

		hostname, err := getBaseDomain(url)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to parse URL: %v", err), http.StatusInternalServerError)
			return
		}

		internicData, err := fetchFromInternic(hostname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		myAPIData, err := fetchFromMyAPI(hostname)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching data from your API: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"internicData": internicData,
			"myAPIData":    myAPIData,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
