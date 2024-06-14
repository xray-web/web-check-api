package checks

import (
	"net/http"
	"time"

	"github.com/xray-web/web-check-api/checks/clients/ip"
	"github.com/xray-web/web-check-api/checks/store/legacyrank"
)

type Checks struct {
	BlockList  *BlockList
	Carbon     *Carbon
	IpAddress  *Ip
	LegacyRank *LegacyRank
	Rank       *Rank
	SocialTags *SocialTags
	Tls        *Tls
}

func NewChecks() *Checks {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Checks{
		BlockList:  NewBlockList(&ip.NetDNSLookup{}),
		Carbon:     NewCarbon(client),
		IpAddress:  NewIp(NewNetIp()),
		LegacyRank: NewLegacyRank(legacyrank.NewInMemoryStore()),
		Rank:       NewRank(client),
		SocialTags: NewSocialTags(client),
		Tls:        NewTls(client),
	}
}
