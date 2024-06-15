package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xray-web/web-check-api/config"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("start server", func(t *testing.T) {
		t.Parallel()

		srv := New(config.New())
		srv.routes()
		ts := httptest.NewServer(srv.CORS(srv.mux))
		defer ts.Close()

		// wait up tot 10 seconds for health check to return 200
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(+10*time.Second))
		defer cancel()
		for {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/health", nil)
			assert.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				break
			}
		}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		assert.NoError(t, err)
	})
}
