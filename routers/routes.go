package routers

import (
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
)

var header controllers.HeaderController
var archive controllers.ArchivesController
var cookies controllers.CookiesController
var carbon controllers.CarbonController
var blockLists controllers.BlockListsController
var dnsServer controllers.DnsServerController
var dns controllers.DnsController
var dnssec controllers.DnssecController
var firewall controllers.FirewallController
var getIP controllers.GetIPController
var hsts controllers.HstsController
var httpSecurity controllers.HttpSecurityController
var legacyRank controllers.LegacyRankController
var getLinks controllers.GetLinksController
var ports controllers.PortsController
var quality controllers.QualityController
var rank controllers.RankController
var redirects controllers.RedirectsController
var socialTags controllers.SocialTagsController
var techStack controllers.TechStackController

func WebCheckRoutes(route *gin.Engine) {
	api := route.Group("/api")
	{
		api.GET("/headers", header.GetHeaders)
		api.GET("/archives", archive.ArchivesHandler)
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
		api.GET("/tech-stack", techStack.TechStackHandler)

	}
}
