package handlers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const hardTimeout = 5 * time.Second

type SitemapIndex struct {
	Sitemaps []Sitemap `xml:"sitemap"`
}

type Sitemap struct {
	Loc string `xml:"loc"`
}

type URLSet struct {
	URLs []URL `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

func HandleSitemap() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			JSONError(w, errors.New("url query parameter is required"), http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlParam, "http://") && !strings.HasPrefix(urlParam, "https://") {
			urlParam = "https://" + urlParam
		}

		sitemapURL := fmt.Sprintf("%s/sitemap.xml", urlParam)

		sitemap, err := fetchSitemap(sitemapURL)
		if err != nil {
			if err == http.ErrHandlerTimeout {
				JSONError(w, fmt.Errorf("request timed-out after %dms", hardTimeout.Milliseconds()), http.StatusRequestTimeout)
				return
			}

			// If sitemap not found, try to fetch it from robots.txt
			if err.Error() == "404" {
				sitemapURL, err = getSitemapURLFromRobotsTxt(urlParam)
				if err != nil {
					JSON(w, map[string]string{"skipped": "No sitemap found"}, http.StatusOK)
					return
				}

				sitemap, err = fetchSitemap(sitemapURL)
				if err != nil {
					JSONError(w, err, http.StatusInternalServerError)
					return
				}
			} else {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		JSON(w, sitemap, http.StatusOK)
	})
}

func fetchSitemap(sitemapURL string) (interface{}, error) {
	client := http.Client{
		Timeout: hardTimeout,
	}

	resp, err := client.Get(sitemapURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("404")
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch sitemap")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sitemap SitemapIndex
	if err := xml.Unmarshal(body, &sitemap); err == nil && len(sitemap.Sitemaps) > 0 {
		return sitemap, nil
	}

	var urlset URLSet
	if err := xml.Unmarshal(body, &urlset); err == nil && len(urlset.URLs) > 0 {
		return urlset, nil
	}

	return nil, errors.New("invalid sitemap format")
}

func getSitemapURLFromRobotsTxt(baseURL string) (string, error) {
	client := http.Client{
		Timeout: hardTimeout,
	}

	resp, err := client.Get(fmt.Sprintf("%s/robots.txt", baseURL))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch robots.txt")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "sitemap:") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", errors.New("no sitemap found in robots.txt")
}
