package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/chromedp/chromedp"
)

type screenshotResponse struct {
	Image string `json:"image"`
}

func takeScreenshot(targetURL string) (*screenshotResponse, error) {
	if targetURL == "" {
		return nil, errors.New("URL is missing from queryStringParameters")
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	parsedURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return nil, errors.New("URL provided is invalid")
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(800, 600),
		chromedp.Navigate(parsedURL.String()),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		return nil, err
	}

	base64Screenshot := base64.StdEncoding.EncodeToString(buf)
	return &screenshotResponse{
		Image: base64Screenshot,
	}, nil
}

func HandleScreenshot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
			return
		}

		data, err := takeScreenshot(targetURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
