package checks

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type HSTSResponse struct {
	Message    string `json:"message"`
	Compatible bool   `json:"compatible"`
	HSTSHeader string `json:"hstsHeader"`
}

type Hsts struct {
	client *http.Client
}

func NewHsts(client *http.Client) *Hsts {
	return &Hsts{client: client}
}

func (h *Hsts) Validate(ctx context.Context, url string) (*HSTSResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	hstsHeader := resp.Header.Get("Strict-Transport-Security")
	if hstsHeader == "" {
		return &HSTSResponse{Message: "Site does not serve any HSTS headers."}, nil
	}

	if !strings.Contains(hstsHeader, "max-age") {
		return &HSTSResponse{Message: "HSTS max-age is less than 10886400."}, nil
	}

	var maxAgeString string
	for _, h := range strings.Split(hstsHeader, " ") {
		if strings.Contains(h, "max-age=") {
			maxAgeString = extractMaxAgeFromHeader(h)
		}
	}

	maxAge, err := strconv.Atoi(maxAgeString)
	if err != nil {
		return nil, err
	}

	if maxAge < 10886400 {
		return &HSTSResponse{Message: "HSTS max-age is less than 10886400."}, nil
	}

	if !strings.Contains(hstsHeader, "includeSubDomains") {
		return &HSTSResponse{Message: "HSTS header does not include all subdomains."}, nil
	}

	if !strings.Contains(hstsHeader, "preload") {
		return &HSTSResponse{Message: "HSTS header does not contain the preload directive."}, nil
	}

	return &HSTSResponse{
		Message:    "Site is compatible with the HSTS preload list!",
		Compatible: true,
		HSTSHeader: hstsHeader,
	}, nil
}

func extractMaxAgeFromHeader(header string) string {
	var maxAge strings.Builder

	for _, b := range header {
		if unicode.IsDigit(b) {
			maxAge.WriteRune(b)
		}
	}

	return maxAge.String()
}
