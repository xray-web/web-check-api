package legacyrank

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var ErrNotFound = errors.New("domain not found")

type Getter interface {
	GetLegacyRank(domain string) (int, error)
}

type GetterFunc func(domain string) (int, error)

func (f GetterFunc) GetLegacyRank(domain string) (int, error) {
	return f(domain)
}

type InMemoryStore struct{}

var once sync.Once
var data map[string]int //map of domain to rank

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

func (s *InMemoryStore) GetLegacyRank(url string) (int, error) {
	once.Do(func() {
		var err error
		data, err = load()
		if err != nil {
			log.Println(err)
		}
	})

	rank, ok := data[url]
	if !ok {
		return -1, ErrNotFound
	}
	return rank, nil
}

func load() (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip", nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	zf, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}
	f, err := zf.Open("top-1m.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	data := make(map[string]int)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		data[record[1]] = rank
	}
	return data, nil
}
