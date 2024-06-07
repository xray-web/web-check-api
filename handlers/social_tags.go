package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type SocialTagsController struct{}

type SocialTags struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	Keywords           string `json:"keywords"`
	CanonicalUrl       string `json:"canonicalUrl"`
	OgTitle            string `json:"ogTitle"`
	OgType             string `json:"ogType"`
	OgImage            string `json:"ogImage"`
	OgUrl              string `json:"ogUrl"`
	OgDescription      string `json:"ogDescription"`
	OgSiteName         string `json:"ogSiteName"`
	TwitterCard        string `json:"twitterCard"`
	TwitterSite        string `json:"twitterSite"`
	TwitterCreator     string `json:"twitterCreator"`
	TwitterTitle       string `json:"twitterTitle"`
	TwitterDescription string `json:"twitterDescription"`
	TwitterImage       string `json:"twitterImage"`
	ThemeColor         string `json:"themeColor"`
	Robots             string `json:"robots"`
	Googlebot          string `json:"googlebot"`
	Generator          string `json:"generator"`
	Viewport           string `json:"viewport"`
	Author             string `json:"author"`
	Publisher          string `json:"publisher"`
	Favicon            string `json:"favicon"`
}

func (ctrl *SocialTagsController) GetSocialTagsHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// Fetch HTML content from the URL
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	// Extract social tags metadata
	tags := &SocialTags{
		Title:              doc.Find("head title").Text(),
		Description:        doc.Find("meta[name='description']").AttrOr("content", ""),
		Keywords:           doc.Find("meta[name='keywords']").AttrOr("content", ""),
		CanonicalUrl:       doc.Find("link[rel='canonical']").AttrOr("href", ""),
		OgTitle:            doc.Find("meta[property='og:title']").AttrOr("content", ""),
		OgType:             doc.Find("meta[property='og:type']").AttrOr("content", ""),
		OgImage:            doc.Find("meta[property='og:image']").AttrOr("content", ""),
		OgUrl:              doc.Find("meta[property='og:url']").AttrOr("content", ""),
		OgDescription:      doc.Find("meta[property='og:description']").AttrOr("content", ""),
		OgSiteName:         doc.Find("meta[property='og:site_name']").AttrOr("content", ""),
		TwitterCard:        doc.Find("meta[name='twitter:card']").AttrOr("content", ""),
		TwitterSite:        doc.Find("meta[name='twitter:site']").AttrOr("content", ""),
		TwitterCreator:     doc.Find("meta[name='twitter:creator']").AttrOr("content", ""),
		TwitterTitle:       doc.Find("meta[name='twitter:title']").AttrOr("content", ""),
		TwitterDescription: doc.Find("meta[name='twitter:description']").AttrOr("content", ""),
		TwitterImage:       doc.Find("meta[name='twitter:image']").AttrOr("content", ""),
		ThemeColor:         doc.Find("meta[name='theme-color']").AttrOr("content", ""),
		Robots:             doc.Find("meta[name='robots']").AttrOr("content", ""),
		Googlebot:          doc.Find("meta[name='googlebot']").AttrOr("content", ""),
		Generator:          doc.Find("meta[name='generator']").AttrOr("content", ""),
		Viewport:           doc.Find("meta[name='viewport']").AttrOr("content", ""),
		Author:             doc.Find("meta[name='author']").AttrOr("content", ""),
		Publisher:          doc.Find("link[rel='publisher']").AttrOr("href", ""),
		Favicon:            doc.Find("link[rel='icon']").AttrOr("href", ""),
	}

	if isEmpty(tags) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no metadata found"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func isEmpty(tags *SocialTags) bool {
	return tags.Title == "" &&
		tags.Description == "" &&
		tags.Keywords == "" &&
		tags.CanonicalUrl == "" &&
		tags.OgTitle == "" &&
		tags.OgType == "" &&
		tags.OgImage == "" &&
		tags.OgUrl == "" &&
		tags.OgDescription == "" &&
		tags.OgSiteName == "" &&
		tags.TwitterCard == "" &&
		tags.TwitterSite == "" &&
		tags.TwitterCreator == "" &&
		tags.TwitterTitle == "" &&
		tags.TwitterDescription == "" &&
		tags.TwitterImage == "" &&
		tags.ThemeColor == "" &&
		tags.Robots == "" &&
		tags.Googlebot == "" &&
		tags.Generator == "" &&
		tags.Viewport == "" &&
		tags.Author == "" &&
		tags.Publisher == "" &&
		tags.Favicon == ""
}

func HandleGetSocialTags() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}

		// Fetch HTML content from the URL
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Parse HTML document
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return
		}

		// Extract social tags metadata
		tags := &SocialTags{
			Title:              doc.Find("head title").Text(),
			Description:        doc.Find("meta[name='description']").AttrOr("content", ""),
			Keywords:           doc.Find("meta[name='keywords']").AttrOr("content", ""),
			CanonicalUrl:       doc.Find("link[rel='canonical']").AttrOr("href", ""),
			OgTitle:            doc.Find("meta[property='og:title']").AttrOr("content", ""),
			OgType:             doc.Find("meta[property='og:type']").AttrOr("content", ""),
			OgImage:            doc.Find("meta[property='og:image']").AttrOr("content", ""),
			OgUrl:              doc.Find("meta[property='og:url']").AttrOr("content", ""),
			OgDescription:      doc.Find("meta[property='og:description']").AttrOr("content", ""),
			OgSiteName:         doc.Find("meta[property='og:site_name']").AttrOr("content", ""),
			TwitterCard:        doc.Find("meta[name='twitter:card']").AttrOr("content", ""),
			TwitterSite:        doc.Find("meta[name='twitter:site']").AttrOr("content", ""),
			TwitterCreator:     doc.Find("meta[name='twitter:creator']").AttrOr("content", ""),
			TwitterTitle:       doc.Find("meta[name='twitter:title']").AttrOr("content", ""),
			TwitterDescription: doc.Find("meta[name='twitter:description']").AttrOr("content", ""),
			TwitterImage:       doc.Find("meta[name='twitter:image']").AttrOr("content", ""),
			ThemeColor:         doc.Find("meta[name='theme-color']").AttrOr("content", ""),
			Robots:             doc.Find("meta[name='robots']").AttrOr("content", ""),
			Googlebot:          doc.Find("meta[name='googlebot']").AttrOr("content", ""),
			Generator:          doc.Find("meta[name='generator']").AttrOr("content", ""),
			Viewport:           doc.Find("meta[name='viewport']").AttrOr("content", ""),
			Author:             doc.Find("meta[name='author']").AttrOr("content", ""),
			Publisher:          doc.Find("link[rel='publisher']").AttrOr("href", ""),
			Favicon:            doc.Find("link[rel='icon']").AttrOr("href", ""),
		}

		if isEmpty(tags) {
			JSONError(w, errors.New("no metadata found"), http.StatusBadRequest)
			return
		}

		JSON(w, tags, http.StatusOK)
	})
}
