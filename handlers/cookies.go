package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Function to get cookies using chromedp
func getChromedpCookies(url string) ([]map[string]interface{}, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Create a timeout context for chromedp actions
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // Increased timeout
	defer cancel()

	var cookies []*network.Cookie
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookieParams := network.GetCookies().WithUrls([]string{url})
			var err error
			cookies, err = cookieParams.Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, err
	}

	var cookiesList []map[string]interface{}
	for _, c := range cookies {
		cookie := map[string]interface{}{
			"name":    c.Name,
			"value":   c.Value,
			"domain":  c.Domain,
			"path":    c.Path,
			"expires": c.Expires,
			// "size":         cdp.CookieSize(c),
			"httpOnly": c.HTTPOnly,
			"secure":   c.Secure,
			"session":  c.Session,
			"sameSite": c.SameSite.String(),
			"priority": c.Priority.String(),
			// "sameParty":    c.SameParty,
			"sourceScheme": c.SourceScheme.String(),
		}
		cookiesList = append(cookiesList, cookie)
	}

	return cookiesList, nil
}

func HandleCookies() http.Handler {
	type Response struct {
		HeaderCookies []string         `json:"headerCookies"`
		ClientCookies []map[string]any `json:"clientCookies"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		var headerCookies []string
		var clientCookies []map[string]interface{}

		url := rawURL.String()
		// Fetch headers using http.Get
		resp, err := http.Get(url)
		if err != nil {
			JSONError(w, fmt.Errorf("request failed: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		headerCookies = resp.Header["Set-Cookie"]

		// Fetch client cookies using chromedp
		clientCookies, err = getChromedpCookies(url)
		if err != nil {
			clientCookies = nil
		}

		if len(headerCookies) == 0 && (len(clientCookies) == 0) {
			JSON(w, KV{"skipped": "No cookies"}, http.StatusOK)
			return
		}
		JSON(w, Response{
			HeaderCookies: headerCookies,
			ClientCookies: clientCookies,
		}, http.StatusOK)

	})
}
