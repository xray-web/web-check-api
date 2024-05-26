package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type FirewallController struct{}

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
	fullURL := ""
	if !strings.HasPrefix(url, "http") {
		fullURL = "http://" + url
	} else {
		fullURL = url
	}

	resp, err := http.Get(fullURL)
	if err != nil {
		return wafResponse{}, fmt.Errorf("error fetching URL: %s", err.Error())
	}
	defer resp.Body.Close()

	headers := resp.Header

	for header, value := range headers {
		if strings.Contains(header, "server") && strings.Contains(value[0], "cloudflare") {
			return wafResponse{HasWaf: true, Waf: cloudflare}, nil
		}

		if strings.Contains(header, "x-powered-by") && strings.Contains(value[0], "AWS Lambda") {
			return wafResponse{HasWaf: true, Waf: awsWAF}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "AkamaiGHost") {
			return wafResponse{HasWaf: true, Waf: akamai}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "Sucuri") {
			return wafResponse{HasWaf: true, Waf: sucuri}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "BarracudaWAF") {
			return wafResponse{HasWaf: true, Waf: barracuda}, nil
		}

		if strings.Contains(header, "server") && (strings.Contains(value[0], "F5 BIG-IP") || strings.Contains(value[0], "BIG-IP")) {
			return wafResponse{HasWaf: true, Waf: f5}, nil
		}

		_, sucuriId := headers[http.CanonicalHeaderKey("x-sucuri-id")]
		_, sucuriCache := headers[http.CanonicalHeaderKey("x-sucuri-cache")]

		if sucuriId || sucuriCache {
			return wafResponse{HasWaf: true, Waf: sucuriProxy}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "FortiWeb") {
			return wafResponse{HasWaf: true, Waf: fortinet}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "Imperva") {
			return wafResponse{HasWaf: true, Waf: imperva}, nil
		}

		if _, exists := headers[http.CanonicalHeaderKey("x-protected-by")]; exists && strings.Contains(headers[http.CanonicalHeaderKey("x-protected-by")][0], "Sqreen") {
			return wafResponse{HasWaf: true, Waf: sqreen}, nil
		}

		if _, exists := headers[http.CanonicalHeaderKey("x-waf-event-info")]; exists {
			return wafResponse{HasWaf: true, Waf: reblaze}, nil
		}

		if _, exists := headers[http.CanonicalHeaderKey("set-cookie")]; exists && strings.Contains(headers[http.CanonicalHeaderKey("set-cookie")][0], "_citrix_ns_id") {
			return wafResponse{HasWaf: true, Waf: citrix}, nil
		}

		_, deniedReason := headers[http.CanonicalHeaderKey("x-denied-reason")]
		_, requestedMethod := headers[http.CanonicalHeaderKey("x-wzws-requested-method")]

		if deniedReason || requestedMethod {
			return wafResponse{HasWaf: true, Waf: wzb}, nil
		}

		if _, exists := headers["x-webcoment"]; exists {
			return wafResponse{HasWaf: true, Waf: webcoment}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "Yundun") {
			return wafResponse{HasWaf: true, Waf: yundun}, nil
		}

		_, wafInfo := headers["x-yd-waf-info"]
		_, ydInfo := headers["x-yd-info"]

		if wafInfo || ydInfo {
			return wafResponse{HasWaf: true, Waf: yundun}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "Safe3WAF") {
			return wafResponse{HasWaf: true, Waf: safe3}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "NAXSI") {
			return wafResponse{HasWaf: true, Waf: naxsi}, nil
		}

		if _, exists := headers["x-datapower-transactionid"]; exists {
			return wafResponse{HasWaf: true, Waf: ibm}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "QRATOR") {
			return wafResponse{HasWaf: true, Waf: qrator}, nil
		}

		if strings.Contains(header, "server") && strings.Contains(value[0], "ddos-guard") {
			return wafResponse{HasWaf: true, Waf: ddosGuard}, nil
		}
	}

	return wafResponse{HasWaf: false}, nil
}

func (ctrl *FirewallController) FirewallHandler(c *gin.Context) {
	domain := c.Query("url")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	result, err := checkWAF(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
