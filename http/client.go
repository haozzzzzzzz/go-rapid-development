package http

import (
	"context"
	"net"
	"net/http"
	"time"
)

var (
	ShortTimeoutRequestClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     5 * time.Minute,
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

	// Deprecated: use ShortTimeoutRequestClient, which is more semantic
	RequestClient = ShortTimeoutRequestClient

	LongTimeoutRequestClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     5 * time.Minute,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{
					Timeout:   20 * time.Second,
					KeepAlive: 10 * time.Minute,
				}

				return d.DialContext(ctx, network, addr)
			},
		},

		Timeout: 20 * time.Second,
	}

	NoTimeoutRequestClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     5 * time.Minute,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{
					Timeout:   20 * time.Second,
					KeepAlive: 10 * time.Minute,
				}

				return d.DialContext(ctx, network, addr)
			},
		},
	}
)
