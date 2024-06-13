package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

const (
	googleSafeBrowsingAPI = "https://safebrowsing.googleapis.com/v4/threatMatches:find"
	urlHausAPI            = "https://urlhaus-api.abuse.ch/v1/host/"
	phishTankAPI          = "https://checkurl.phishtank.com/checkurl/"
	cloudmersiveAPI       = "https://api.cloudmersive.com/virus/scan/website"
)

type GoogleSafeBrowsingResponse struct {
	Matches []Match `json:"matches,omitempty"`
}

type Match struct {
	ThreatType          string   `json:"threatType,omitempty"`
	PlatformType        string   `json:"platformType,omitempty"`
	ThreatEntryType     string   `json:"threatEntryType,omitempty"`
	Threat              Threat   `json:"threat,omitempty"`
	CacheDuration       string   `json:"cacheDuration,omitempty"`
	HashPrefix          string   `json:"hashPrefix,omitempty"`
	ThreatEntryMetadata Metadata `json:"threatEntryMetadata,omitempty"`
}

type Threat struct {
	URL string `json:"url,omitempty"`
}

type Metadata struct {
	Entries []Entry `json:"entries,omitempty"`
}

type Entry struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type UrlHausResponse struct {
	URLs []string `json:"urls,omitempty"`
}

type PhishTankResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CloudmersiveResponse struct {
	Error     string `json:"error,omitempty"`
	IsSuccess bool   `json:"isSuccess,omitempty"`
	Message   string `json:"message,omitempty"`
}

func getGoogleSafeBrowsingResult(urlParam string) (*GoogleSafeBrowsingResponse, error) {
	apiKey := "" // Set your Google Safe Browsing API key here
	if apiKey == "" {
		return nil, errors.New("GOOGLE_CLOUD_API_KEY is required for the Google Safe Browsing check")
	}

	requestBody := map[string]interface{}{
		"threatInfo": map[string]interface{}{
			"threatTypes":      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION", "API_ABUSE"},
			"platformTypes":    []string{"ANY_PLATFORM"},
			"threatEntryTypes": []string{"URL"},
			"threatEntries": []map[string]string{
				{"url": urlParam},
			},
		},
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(fmt.Sprintf("%s?key=%s", googleSafeBrowsingAPI, apiKey), "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var googleSafeBrowsingResponse GoogleSafeBrowsingResponse
	if err := json.NewDecoder(response.Body).Decode(&googleSafeBrowsingResponse); err != nil {
		return nil, err
	}

	return &googleSafeBrowsingResponse, nil
}

func getUrlHausResult(urlParam string) (*UrlHausResponse, error) {
	parsedURL, err := url.Parse(urlParam)
	if err != nil {
		return nil, err
	}

	domain := parsedURL.Hostname()
	if domain == "" {
		return nil, errors.New("invalid URL format")
	}

	response, err := http.PostForm(urlHausAPI, url.Values{"host": {domain}})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var urlHausResponse UrlHausResponse
	if err := json.Unmarshal(body, &urlHausResponse); err != nil {
		return nil, err
	}

	return &urlHausResponse, nil
}

func getPhishTankResult(urlParam string) (*PhishTankResponse, error) {
	encodedURL := base64.StdEncoding.EncodeToString([]byte(urlParam))
	endpoint := fmt.Sprintf("%s?url=%s", phishTankAPI, encodedURL)

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// Extract title and body
	title := extractTitle(doc)
	bodyContent := extractBody(doc)

	// Create PhishTankResponse
	phishTankResponse := &PhishTankResponse{
		Title: title,
		Body:  bodyContent,
	}

	return phishTankResponse, nil
}

func extractTitle(doc *html.Node) string {
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = strings.TrimSpace(n.FirstChild.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}

func extractBody(doc *html.Node) string {
	var bodyContent string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					bodyContent += strings.TrimSpace(c.Data) + "\n"
				}
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return bodyContent
}

func getCloudmersiveResult(urlParam string) (*CloudmersiveResponse, error) {
	apiKey := "" // Set your Cloudmersive API key here
	if apiKey == "" {
		return nil, errors.New("CLOUDMERSIVE_API_KEY is required for the Cloudmersive check")
	}

	data := url.Values{}
	data.Set("Url", urlParam)

	client := &http.Client{}
	request, err := http.NewRequest("POST", cloudmersiveAPI, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Apikey", apiKey)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var cloudmersiveResponse CloudmersiveResponse
	if err := json.NewDecoder(response.Body).Decode(&cloudmersiveResponse); err != nil {
		return nil, err
	}

	return &cloudmersiveResponse, nil
}

func HandleThreats() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		urlHausResult, err := getUrlHausResult(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		phishTankResult, err := getPhishTankResult(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		cloudmersiveResult, err := getCloudmersiveResult(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		googleSafeBrowsingResult, err := getGoogleSafeBrowsingResult(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"urlHaus":            urlHausResult,
			"phishTank":          phishTankResult,
			"cloudmersive":       cloudmersiveResult,
			"googleSafeBrowsing": googleSafeBrowsingResult,
		}

		JSON(w, response, http.StatusOK)
	})
}
