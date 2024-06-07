package routers

import (
	"github.com/xray-web/web-check-api/handlers"

	"github.com/gin-gonic/gin"
)

var header handlers.HeaderController
var cookies handlers.CookiesController
var carbon handlers.CarbonController
var blockLists handlers.BlockListsController
var dnsServer handlers.DnsServerController
var dns handlers.DnsController
var dnssec handlers.DnssecController
var firewall handlers.FirewallController
var getIP handlers.GetIPController
var hsts handlers.HstsController
var httpSecurity handlers.HttpSecurityController
var legacyRank handlers.LegacyRankController
var getLinks handlers.GetLinksController
var ports handlers.PortsController
var quality handlers.QualityController
var rank handlers.RankController
var redirects handlers.RedirectsController
var socialTags handlers.SocialTagsController
var tls handlers.TlsController
var trace handlers.TraceRouteController

func WebCheckRoutes(route *gin.Engine) {
	api := route.Group("/api")
	{
		api.GET("/headers", header.GetHeaders)
		api.GET("/cookies", cookies.CookiesHandler)
		api.GET("/carbon", carbon.CarbonHandler)
		api.GET("/block-lists", blockLists.BlockListsHandler)
		api.GET("/dns-server", dnsServer.DnsServerHandler)
		api.GET("/dns", dns.DnsHandler)
		api.GET("/dnssec", dnssec.DnssecHandler)
		api.GET("/firewall", firewall.FirewallHandler)
		api.GET("/get-ip", getIP.GetIPHandler)
		api.GET("/hsts", hsts.HstsHandler)
		api.GET("/http-security", httpSecurity.HttpSecurityHandler)
		api.GET("/legacy-rank", legacyRank.LegacyRankHandler)
		api.GET("/linked-pages", getLinks.GetLinksHandler)
		api.GET("/ports", ports.GetPortsHandler)
		api.GET("/quality", quality.GetQualityHandler)
		api.GET("/rank", rank.GetRankHandler)
		api.GET("/redirects", redirects.GetRedirectsHandler)
		api.GET("/social-tags", socialTags.GetSocialTagsHandler)
		api.GET("/tls", tls.TlsHandler)
		api.GET("/trace-route", trace.TracerouteHandler)

	}
}
