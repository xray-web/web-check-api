package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TrancoRanks struct {
	Ranks []TrancoRank `json:"ranks"`
}

type TrancoRank struct {
	Date string `json:"date"`
	Rank int    `json:"rank"`
}

type Rank struct {
	client *http.Client
}

func NewRank(client *http.Client) *Rank {
	return &Rank{client: client}
}

func (r *Rank) GetRank(ctx context.Context, url string) (*TrancoRanks, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://tranco-list.eu/api/ranks/domain/%s", url), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TrancoRanks
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}
