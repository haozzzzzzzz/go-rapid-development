package http

import (
	"context"
	"net"
	"net/http"
	"time"
)

var (
	RequestClient = &http.Client{
		Transport: &http.Transport{
			//MaxIdleConns:        1000,
			//MaxIdleConnsPerHost: 200,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{
					Timeout:   2 * time.Second,
					KeepAlive: 10 * time.Minute,
				}

				return d.DialContext(ctx, network, addr)
			},
		},

		Timeout: 2 * time.Second,
	}
)
