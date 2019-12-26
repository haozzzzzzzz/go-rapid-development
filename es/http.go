package es

import (
	"net"
	"net/http"
	"time"
)

var ShortTimeoutTransport = &http.Transport{
	MaxIdleConnsPerHost:   100,
	ResponseHeaderTimeout: 2 * time.Second,
	IdleConnTimeout:       1 * time.Minute,
	DialContext: (&net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 10 * time.Minute,
	}).DialContext,
}

var TimeoutTransport = ShortTimeoutTransport

var NoTimeoutTransport = &http.Transport{
	MaxIdleConns:    1000,
	MaxConnsPerHost: 100,
	IdleConnTimeout: 1 * time.Minute,
	DialContext: (&net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 10 * time.Minute,
	}).DialContext,
}
