package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var SECURITY_TXT_PATHS = []string{
	"/security.txt",
	"/.well-known/security.txt",
}

func parseResult(result string) map[string]string {
	output := make(map[string]string)
	counts := make(map[string]int)
	lines := strings.Split(result, "\n")
	regex := ":\\s*"

	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "-----") && strings.TrimSpace(line) != "" {
			parts := strings.SplitN(line, regex, 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if _, exists := output[key]; exists {
					counts[key]++
					key = fmt.Sprintf("%s%d", key, counts[key])
				}
				output[key] = value
			}
		}
	}

	return output
}

func isPgpSigned(result string) bool {
	return strings.Contains(result, "-----BEGIN PGP SIGNED MESSAGE-----")
}

func HandleSecurityTxt() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		rawURL.Path = ""
		for _, path := range SECURITY_TXT_PATHS {
			result, err := fetchSecurityTxt(rawURL, path)
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}

			if result != "" && strings.Contains(result, "<html") {
				JSON(w, map[string]bool{"isPresent": false}, http.StatusOK)
				return
			}

			if result != "" {
				response := map[string]interface{}{
					"isPresent":   true,
					"foundIn":     path,
					"content":     result,
					"isPgpSigned": isPgpSigned(result),
					"fields":      parseResult(result),
				}
				JSON(w, response, http.StatusOK)
				return
			}
		}

		JSON(w, map[string]bool{"isPresent": false}, http.StatusOK)
	})
}

func fetchSecurityTxt(baseURL *url.URL, path string) (string, error) {
	secTxtURL, err := baseURL.Parse(path)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(secTxtURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
