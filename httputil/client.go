package httputil

import (
	"context"
	"net"
	"net/http"
	"time"
)

func NewClient(network, addr string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		},
		Timeout: time.Second * 10,
	}
}
