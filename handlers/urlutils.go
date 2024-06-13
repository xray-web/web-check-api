package handlers

import (
	"net/http"
	"net/url"
)

func extractURL(r *http.Request) (*url.URL, error) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		return nil, ErrMissingURLParameter
	}
	u, err := url.Parse(rawURL) // this will parse almost anything so hard to trigger, we do a mosre stric check after cleaning up the URL
	if err != nil {
		return nil, ErrInvalidURL
	}

	// If the url has no scheme then its a relative path so we can add the scheme and parse again
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	u, err = url.Parse(u.String())
	if err != nil {
		return nil, ErrInvalidURL // to ensure u is nil
	}
	// with a fixed up url we cand do a more strict parse
	u, err = url.ParseRequestURI(u.String())
	if err != nil {
		return nil, ErrInvalidURL
	}
	return u, nil
}
