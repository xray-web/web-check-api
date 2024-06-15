package checks

import (
	"context"
	"net/http"
)

type Headers struct {
	client *http.Client
}

func NewHeaders(client *http.Client) *Headers {
	return &Headers{client: client}
}

func (h *Headers) List(ctx context.Context, url string) (http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Header, nil
}
