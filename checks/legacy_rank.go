package checks

import "github.com/xray-web/web-check-api/checks/store/legacyrank"

type DomainRank struct {
	Domain string `json:"domain"`
	Rank   int    `json:"rank"`
}

type LegacyRank struct {
	data legacyrank.Getter
}

func NewLegacyRank(lrg legacyrank.Getter) *LegacyRank {
	return &LegacyRank{data: lrg}
}

func (lr *LegacyRank) LegacyRank(domain string) (*DomainRank, error) {
	rank, err := lr.data.GetLegacyRank(domain)
	if err != nil {
		return nil, err
	}
	return &DomainRank{
		Domain: domain,
		Rank:   rank,
	}, nil
}
