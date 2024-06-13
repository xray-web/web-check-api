package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/chromedp/chromedp"
)

type screenshotResponse struct {
	Image string `json:"image"`
}

func takeScreenshot(targetURL string) (*screenshotResponse, error) {
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
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		data, err := takeScreenshot(rawURL.String())
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
