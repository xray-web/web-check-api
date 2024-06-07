package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type GetLinksController struct{}

type LinkResponse struct {
	Internal []string `json:"internal"`
	External []string `json:"external"`
}

type ErrorResponse struct {
	Skipped string `json:"skipped"`
}

func (e ErrorResponse) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (ctrl *GetLinksController) GetLinksHandler(c *gin.Context) {
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Ensure the URL has a scheme
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	internalLinks, externalLinks, err := getLinks(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(internalLinks) == 0 && len(externalLinks) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Skipped: "No internal or external links found. " +
				"This may be due to the website being dynamically rendered, using a client-side framework (like React), and without SSR enabled. " +
				"That would mean that the static HTML returned from the HTTP request doesn't contain any meaningful content for Web-Check to analyze. " +
				"You can rectify this by using a headless browser to render the page instead.",
		})
		return
	}

	c.JSON(http.StatusOK, LinkResponse{
		Internal: internalLinks,
		External: externalLinks,
	})
}

func getLinks(targetURL string) ([]string, []string, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("received non-200 response code")
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %s", err)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing HTML: %s", err)
	}

	internalLinksMap := make(map[string]int)
	externalLinksMap := make(map[string]int)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					absoluteURL := resolveURL(parsedURL, href)
					if isInternalLink(absoluteURL, parsedURL) {
						internalLinksMap[absoluteURL]++
					} else if strings.HasPrefix(absoluteURL, "http://") || strings.HasPrefix(absoluteURL, "https://") {
						externalLinksMap[absoluteURL]++
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)

	internalLinks := sortAndExtractKeys(internalLinksMap)
	externalLinks := sortAndExtractKeys(externalLinksMap)

	return internalLinks, externalLinks, nil
}

func resolveURL(baseURL *url.URL, href string) string {
	u, err := url.Parse(href)
	if err != nil || !u.IsAbs() {
		return baseURL.ResolveReference(u).String()
	}
	return u.String()
}

func isInternalLink(link string, baseURL *url.URL) bool {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return false
	}
	return parsedLink.Hostname() == baseURL.Hostname()
}

func sortAndExtractKeys(m map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var sortedKeys []string
	for _, kv := range ss {
		sortedKeys = append(sortedKeys, kv.Key)
	}

	return sortedKeys
}

func HandleGetLinks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		// Ensure the URL has a scheme
		if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
			targetURL = "http://" + targetURL
		}

		internalLinks, externalLinks, err := getLinks(targetURL)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if len(internalLinks) == 0 && len(externalLinks) == 0 {
			JSONError(w, ErrorResponse{
				Skipped: "No internal or external links found. " +
					"This may be due to the website being dynamically rendered, using a client-side framework (like React), and without SSR enabled. " +
					"That would mean that the static HTML returned from the HTTP request doesn't contain any meaningful content for Web-Check to analyze. " +
					"You can rectify this by using a headless browser to render the page instead.",
			}, http.StatusBadRequest)
			return
		}

		JSON(w, LinkResponse{
			Internal: internalLinks,
			External: externalLinks,
		}, http.StatusOK)
	})
}
