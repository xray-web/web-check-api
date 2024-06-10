package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	cloudflare  = "Cloudflare"
	awsWAF      = "AWS WAF"
	akamai      = "Akamai"
	sucuri      = "Sucuri"
	barracuda   = "Barracuda WAF"
	f5          = "F5 BIG-IP"
	sucuriProxy = "Sucuri CloudProxy WAF"
	fortinet    = "Fortinet FortiWeb WAF"
	imperva     = "Imperva SecureSphere WAF"
	sqreen      = "Sqreen"
	reblaze     = "Reblaze WAF"
	citrix      = "Citrix NetScaler"
	wzb         = "WangZhanBao WAF"
	webcoment   = "Webcoment Firewall"
	yundun      = "Yundun WAF"
	safe3       = "Safe3 Web Application Firewall"
	naxsi       = "NAXSI WAF"
	ibm         = "IBM WebSphere DataPower"
	qrator      = "QRATOR WAF"
	ddosGuard   = "DDoS-Guard WAF"
)

type wafResponse struct {
	HasWaf bool   `json:"hasWaf"`
	Waf    string `json:"waf,omitempty"`
}

func checkWAF(url string) (wafResponse, error) {
	// TODO(Lissy93): does this test require we set scheme to http?
	resp, err := http.Get(url)
	if err != nil {
		return wafResponse{}, fmt.Errorf("error fetching URL: %s", err.Error())
	}
	defer resp.Body.Close()

	headers := resp.Header

	for header, values := range headers {
		lowerHeader := strings.ToLower(header)

		for _, value := range values {
			lowerValue := strings.ToLower(value)

			switch {
			case lowerHeader == "server" && strings.Contains(lowerValue, "cloudflare"):
				return wafResponse{HasWaf: true, Waf: cloudflare}, nil
			case lowerHeader == "x-powered-by" && strings.Contains(lowerValue, "aws lambda"):
				return wafResponse{HasWaf: true, Waf: awsWAF}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "akamaighost"):
				return wafResponse{HasWaf: true, Waf: akamai}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "sucuri"):
				return wafResponse{HasWaf: true, Waf: sucuri}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "barracudawaf"):
				return wafResponse{HasWaf: true, Waf: barracuda}, nil
			case lowerHeader == "server" && (strings.Contains(lowerValue, "f5 big-ip") || strings.Contains(lowerValue, "big-ip")):
				return wafResponse{HasWaf: true, Waf: f5}, nil
			case lowerHeader == "x-sucuri-id" || lowerHeader == "x-sucuri-cache":
				return wafResponse{HasWaf: true, Waf: sucuriProxy}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "fortiweb"):
				return wafResponse{HasWaf: true, Waf: fortinet}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "imperva"):
				return wafResponse{HasWaf: true, Waf: imperva}, nil
			case lowerHeader == "x-protected-by" && strings.Contains(lowerValue, "sqreen"):
				return wafResponse{HasWaf: true, Waf: sqreen}, nil
			case lowerHeader == "x-waf-event-info":
				return wafResponse{HasWaf: true, Waf: reblaze}, nil
			case lowerHeader == "set-cookie" && strings.Contains(lowerValue, "_citrix_ns_id"):
				return wafResponse{HasWaf: true, Waf: citrix}, nil
			case lowerHeader == "x-denied-reason" || lowerHeader == "x-wzws-requested-method":
				return wafResponse{HasWaf: true, Waf: wzb}, nil
			case lowerHeader == "x-webcoment":
				return wafResponse{HasWaf: true, Waf: webcoment}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "yundun"):
				return wafResponse{HasWaf: true, Waf: yundun}, nil
			case lowerHeader == "x-yd-waf-info" || lowerHeader == "x-yd-info":
				return wafResponse{HasWaf: true, Waf: yundun}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "safe3waf"):
				return wafResponse{HasWaf: true, Waf: safe3}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "naxsi"):
				return wafResponse{HasWaf: true, Waf: naxsi}, nil
			case lowerHeader == "x-datapower-transactionid":
				return wafResponse{HasWaf: true, Waf: ibm}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "qrator"):
				return wafResponse{HasWaf: true, Waf: qrator}, nil
			case lowerHeader == "server" && strings.Contains(lowerValue, "ddos-guard"):
				return wafResponse{HasWaf: true, Waf: ddosGuard}, nil
			}
		}
	}

	return wafResponse{HasWaf: false}, nil
}

func HandleFirewall() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := checkWAF(rawURL.String())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
