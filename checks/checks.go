package checks

import (
	"net/http"
	"time"
)

type Checks struct {
	Carbon *Carbon
	Rank   *Rank
}

func NewChecks() *Checks {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Checks{
		Carbon: NewCarbon(client),
		Rank:   NewRank(client),
	}
}
