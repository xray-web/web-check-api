package checks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type LinkedPagesData struct {
	Internal []string `json:"internal"`
	External []string `json:"external"`
}

type LinkedPages struct {
	client *http.Client
}

func NewLinkedPages(client *http.Client) *LinkedPages {
	return &LinkedPages{client: client}
}

func (l *LinkedPages) GetLinkedPages(ctx context.Context, targetURL *url.URL) (LinkedPagesData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
	if err != nil {
		return LinkedPagesData{}, err
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return LinkedPagesData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LinkedPagesData{}, fmt.Errorf("received non-200 response code")
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return LinkedPagesData{}, err
	}

	internalLinksMap := make(map[string]int)
	externalLinksMap := make(map[string]int)
	walkDom(doc, targetURL, internalLinksMap, externalLinksMap)

	return LinkedPagesData{
		Internal: sortURLsByFrequency(internalLinksMap),
		External: sortURLsByFrequency(externalLinksMap),
	}, nil
}

func walkDom(n *html.Node, parsedTargetURL *url.URL, internalLinksMap map[string]int, externalLinksMap map[string]int) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				href := attr.Val
				absoluteURL, err := resolveURL(parsedTargetURL, href)
				if err != nil {
					continue
				}
				if strings.TrimPrefix(absoluteURL.Hostname(), "www.") == parsedTargetURL.Hostname() {
					internalLinksMap[absoluteURL.String()]++
				} else if absoluteURL.Scheme == "http" || absoluteURL.Scheme == "https" {
					externalLinksMap[absoluteURL.String()]++
				}
				break
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walkDom(child, parsedTargetURL, internalLinksMap, externalLinksMap)
	}
}

// NOTE: This function resolves a href based on how it would be interpreted by the browser, and does NOT check for typos or whether the URL is reachable.
// Only hrefs containing a scheme or beginning with "//" (denoting a relative scheme) will be resolved as absolute URLs.
// E.g. A href of "http//example.com" will resolve against a base url of "http://example.com" as "http://example.com/http//example.com" since this is how the browser will interpret it
func resolveURL(baseURL *url.URL, href string) (*url.URL, error) {
	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		return baseURL.ResolveReference(u), nil
	}
	return u, nil
}

func sortURLsByFrequency(linksMap map[string]int) []string {
	type link struct {
		URL       string
		Frequency int
	}

	var links []link
	for k, v := range linksMap {
		links = append(links, link{k, v})
	}

	sort.SliceStable(links, func(i, j int) bool {
		return links[i].URL < links[j].URL
	})

	sort.SliceStable(links, func(i, j int) bool {
		return links[i].Frequency > links[j].Frequency
	})

	var sortedLinks []string
	for _, link := range links {
		sortedLinks = append(sortedLinks, link.URL)
	}

	return sortedLinks
}
