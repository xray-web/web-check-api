package checks

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

const (
	StrictTransportSecurity = "Strict-Transport-Security"
	includeSubDomains       = "includeSubDomains"
	preload                 = "preload"
	NilHeadersError         = "Site does not serve any HSTS headers."
	MaxAgeError             = "HSTS max-age is less than 10886400."
	SubdomainsError         = "HSTS header does not include all subdomains."
	PreloadError            = "HSTS header does not contain the preload directive."
	HstsSuccess             = "Site is compatible with the HSTS preload list!"
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

	hstsHeader := resp.Header.Get(StrictTransportSecurity)
	if hstsHeader == "" {
		return &HSTSResponse{
			Message: NilHeadersError,
		}, nil
	}

	maxAge := extractMaxAgeFromHeader(hstsHeader)
	if maxAge == "" {
		return &HSTSResponse{Message: MaxAgeError}, nil
	}

	maxAgeInt, err := convertMaxAgeStringToInt(maxAge)
	if err != nil {
		return nil, err
	}

	if maxAgeInt < 10886400 {
		return &HSTSResponse{Message: MaxAgeError}, nil
	}

	if !strings.Contains(hstsHeader, includeSubDomains) {
		return &HSTSResponse{Message: SubdomainsError}, nil
	}

	if !strings.Contains(hstsHeader, preload) {
		return &HSTSResponse{Message: PreloadError}, nil
	}

	return &HSTSResponse{
		Message:    HstsSuccess,
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

func convertMaxAgeStringToInt(maxAge string) (int, error) {
	return strconv.Atoi(maxAge)
}
