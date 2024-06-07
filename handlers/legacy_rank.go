package handlers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	fileURL      = "https://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"
	tempFilePath = "/tmp/top-1m.csv"
)

type RankResponse struct {
	Domain  string `json:"domain"`
	Rank    string `json:"rank"`
	IsFound bool   `json:"isFound"`
}

func checkLegacyRank(urlStr string) (RankResponse, error) {
	var domain string
	var err error

	// Parse the URL to extract the domain
	u, err := url.Parse(urlStr)
	if err != nil {
		return RankResponse{}, fmt.Errorf("invalid URL")
	}

	// Extract the domain from the parsed URL
	if u.Host != "" {
		domain = u.Host
	} else {
		// If Host is empty, try to extract the domain from the Path
		parts := strings.Split(u.Path, "/")
		if len(parts) > 0 {
			domain = parts[0]
		} else {
			return RankResponse{}, fmt.Errorf("unable to extract domain from URL")
		}
	}

	// Download and unzip the file if not in cache
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		if err := downloadAndUnzip(fileURL); err != nil {
			return RankResponse{}, err
		}
	}

	// Parse the CSV and find the rank
	file, err := os.Open(tempFilePath)
	if err != nil {
		return RankResponse{}, fmt.Errorf("error opening CSV file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return RankResponse{}, fmt.Errorf("error reading CSV record: %s", err)
		}

		if record[1] == domain {
			return RankResponse{
				Domain:  domain,
				Rank:    record[0],
				IsFound: true,
			}, nil
		}
	}

	return RankResponse{
		Domain:  domain,
		IsFound: false,
	}, nil
}

func downloadAndUnzip(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %s", err)
	}
	defer resp.Body.Close()

	zipFile, err := os.Create(tempFilePath + ".zip")
	if err != nil {
		return fmt.Errorf("error creating zip file: %s", err)
	}
	defer zipFile.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing zip file: %s", err)
	}

	err = unzip(tempFilePath+".zip", "/tmp")
	if err != nil {
		return fmt.Errorf("error unzipping file: %s", err)
	}

	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func HandleLegacyRank() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		result, err := checkLegacyRank(url)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, result, http.StatusOK)
	})
}
