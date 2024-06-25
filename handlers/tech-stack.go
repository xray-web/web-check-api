package handlers

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func extractTechnologies(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	technologies := make([]string, 0)

	doc.Find("script[src], link[href]").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			technologies = append(technologies, src)
		}
		if href, exists := s.Attr("href"); exists {
			technologies = append(technologies, href)
		}
	})

	return technologies, nil
}

func HandleTechStack() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		technologies, err := extractTechnologies(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if len(technologies) == 0 {
			JSONError(w, errors.New("Unable to find any technologies for site"), http.StatusInternalServerError)
			return
		}

		JSON(w, technologies, http.StatusOK)
	})
}
