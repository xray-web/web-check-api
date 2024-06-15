package checks

import (
	"net/http"
	"time"

	"github.com/xray-web/web-check-api/checks/clients/ip"
	"github.com/xray-web/web-check-api/checks/store/legacyrank"
)

type Checks struct {
	BlockList   *BlockList
	Carbon      *Carbon
	Headers     *Headers
	IpAddress   *Ip
	LegacyRank  *LegacyRank
	LinkedPages *LinkedPages
	Rank        *Rank
	SocialTags  *SocialTags
	Tls         *Tls
}

func NewChecks() *Checks {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Checks{
		BlockList:   NewBlockList(&ip.NetDNSLookup{}),
		Carbon:      NewCarbon(client),
		Headers:     NewHeaders(client),
		IpAddress:   NewIp(NewNetIp()),
		LegacyRank:  NewLegacyRank(legacyrank.NewInMemoryStore()),
		LinkedPages: NewLinkedPages(client),
		Rank:        NewRank(client),
		SocialTags:  NewSocialTags(client),
		Tls:         NewTls(client),
	}
}
