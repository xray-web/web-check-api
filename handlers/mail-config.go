package handlers

import (
	"net/http"
	"strings"

	"github.com/miekg/dns"
)

func ResolveMx(domain string) ([]*dns.MX, int, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return nil, dns.RcodeServerFailure, err
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, r.Rcode, &dns.Error{}
	}
	var mxRecords []*dns.MX
	for _, ans := range r.Answer {
		if mx, ok := ans.(*dns.MX); ok {
			mxRecords = append(mxRecords, mx)
		}
	}
	if len(mxRecords) == 0 {
		return nil, dns.RcodeNameError, nil
	}
	return mxRecords, dns.RcodeSuccess, nil
}

func ResolveTxt(domain string) ([]string, int, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)
	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return nil, dns.RcodeServerFailure, err
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, r.Rcode, &dns.Error{}
	}
	var txtRecords []string
	for _, ans := range r.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			txtRecords = append(txtRecords, txt.Txt...)
		}
	}
	if len(txtRecords) == 0 {
		return nil, dns.RcodeNameError, nil
	}
	return txtRecords, dns.RcodeSuccess, nil
}

func HandleMailConfig() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawURL, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		mxRecords, rcode, err := ResolveMx(rawURL.Hostname())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if rcode == dns.RcodeNameError || rcode == dns.RcodeServerFailure {
			JSON(w, map[string]string{"skipped": "No mail server in use on this domain"}, http.StatusOK)
			return
		}

		txtRecords, rcode, err := ResolveTxt(rawURL.Hostname())
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if rcode == dns.RcodeNameError || rcode == dns.RcodeServerFailure {
			JSON(w, map[string]string{"skipped": "No mail server in use on this domain"}, http.StatusOK)
			return
		}

		emailTxtRecords := filterEmailTxtRecords(txtRecords)
		mailServices := identifyMailServices(emailTxtRecords, mxRecords)

		JSON(w, map[string]interface{}{
			"mxRecords":    mxRecords,
			"txtRecords":   emailTxtRecords,
			"mailServices": mailServices,
		}, http.StatusOK)
	})
}

func filterEmailTxtRecords(records []string) []string {
	var emailTxtRecords []string
	for _, record := range records {
		if strings.HasPrefix(record, "v=spf1") ||
			strings.HasPrefix(record, "v=DKIM1") ||
			strings.HasPrefix(record, "v=DMARC1") ||
			strings.HasPrefix(record, "protonmail-verification=") ||
			strings.HasPrefix(record, "google-site-verification=") ||
			strings.HasPrefix(record, "MS=") ||
			strings.HasPrefix(record, "zoho-verification=") ||
			strings.HasPrefix(record, "titan-verification=") ||
			strings.Contains(record, "bluehost.com") {
			emailTxtRecords = append(emailTxtRecords, record)
		}
	}
	return emailTxtRecords
}

func identifyMailServices(emailTxtRecords []string, mxRecords []*dns.MX) []map[string]string {
	var mailServices []map[string]string
	for _, record := range emailTxtRecords {
		if strings.HasPrefix(record, "protonmail-verification=") {
			mailServices = append(mailServices, map[string]string{"provider": "ProtonMail", "value": strings.Split(record, "=")[1]})
		} else if strings.HasPrefix(record, "google-site-verification=") {
			mailServices = append(mailServices, map[string]string{"provider": "Google Workspace", "value": strings.Split(record, "=")[1]})
		} else if strings.HasPrefix(record, "MS=") {
			mailServices = append(mailServices, map[string]string{"provider": "Microsoft 365", "value": strings.Split(record, "=")[1]})
		} else if strings.HasPrefix(record, "zoho-verification=") {
			mailServices = append(mailServices, map[string]string{"provider": "Zoho", "value": strings.Split(record, "=")[1]})
		} else if strings.HasPrefix(record, "titan-verification=") {
			mailServices = append(mailServices, map[string]string{"provider": "Titan", "value": strings.Split(record, "=")[1]})
		} else if strings.Contains(record, "bluehost.com") {
			mailServices = append(mailServices, map[string]string{"provider": "BlueHost", "value": record})
		}
	}

	for _, mx := range mxRecords {
		if strings.Contains(mx.Mx, "yahoodns.net") {
			mailServices = append(mailServices, map[string]string{"provider": "Yahoo", "value": mx.Mx})
		} else if strings.Contains(mx.Mx, "mimecast.com") {
			mailServices = append(mailServices, map[string]string{"provider": "Mimecast", "value": mx.Mx})
		}
	}

	return mailServices
}
