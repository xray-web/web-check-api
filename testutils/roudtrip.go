package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func Response(status int, body []byte) *http.Response {
	return &http.Response{
		Status:        http.StatusText(status),
		StatusCode:    status,
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}

func ResponseJSON(status int, body any) *http.Response {
	b, _ := json.Marshal(body)
	return &http.Response{
		Status:        http.StatusText(status),
		StatusCode:    status,
		Body:          io.NopCloser(bytes.NewReader(b)),
		ContentLength: int64(len(b)),
	}
}

func MockClient(responses ...*http.Response) *http.Client {
	return &http.Client{
		Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if len(responses) == 0 {
				return nil, io.EOF
			}
			response := responses[0]
			responses = responses[1:]
			return response, nil
		}),
	}
}
