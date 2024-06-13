package checks

import (
	"context"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type SocialTagsData struct {
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

func (s SocialTagsData) Empty() bool {
	return (SocialTagsData{}) == s
}

type SocialTags struct {
	client *http.Client
}

func NewSocialTags(client *http.Client) *SocialTags {
	return &SocialTags{client: client}
}

func (s *SocialTags) GetSocialTags(ctx context.Context, url string) (*SocialTagsData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract social tags metadata
	tags := &SocialTagsData{
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
	return tags, nil
}
