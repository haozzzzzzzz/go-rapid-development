package es

import (
	"net"
	"net/http"
	"time"
)

var ShortTimeoutTransport = &http.Transport{
	MaxIdleConnsPerHost:   10,
	ResponseHeaderTimeout: time.Second,
	DialContext: (&net.Dialer{
		Timeout: 2 * time.Second,
	}).DialContext,
}

var TimeoutTransport = ShortTimeoutTransport
