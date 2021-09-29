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
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
		Timeout: time.Second * 10,
	}
}
