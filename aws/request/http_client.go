package request

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"time"
)

type HTTPClientSettings struct {
	ConnectTimeout        time.Duration
	ConnKeepAlive         time.Duration
	ExpectContinueTimeout time.Duration
	IdleConnTimeout       time.Duration
	MaxAllIdleConns       int // Zero means no limit
	MaxHostIdleConns      int // If zero, DefaultMaxIdleConnsPerHost is used.
	ResponseHeaderTimeout time.Duration
	TLSHandshakeTimeout   time.Duration

	// Timeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A Timeout of zero means no timeout.
	//
	// The Client cancels requests to the underlying Transport
	// as if the Request's Context ended.
	//
	// For compatibility, the Client will also use the deprecated
	// CancelRequest method on Transport if found. New
	// RoundTripper implementations should use the Request's Context
	// for cancelation instead of implementing CancelRequest.
	Timeout time.Duration
}

func NewHTTPClientWithTimeouts(httpSettings HTTPClientSettings) *http.Client {
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeaderTimeout,
		//Proxy:                 http.ProxyFromEnvironment, // disable environment proxy
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			Timeout:   httpSettings.ConnectTimeout,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConnTimeout,
		TLSHandshakeTimeout:   httpSettings.TLSHandshakeTimeout,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinueTimeout,
	}

	// So client makes HTTP/2 requests
	err := http2.ConfigureTransport(tr)
	if err != nil {
		logrus.Errorf("http2 configure transport failed. %s", err)
	}

	return &http.Client{
		Transport: tr,
		Timeout:   httpSettings.Timeout,
	}
}

var (
	DefaultAWSHttpClientNoTimeout = NewHTTPClientWithTimeouts(HTTPClientSettings{
		ConnectTimeout:        5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ConnKeepAlive:         30 * time.Second,
		MaxAllIdleConns:       10000, // no limit
		MaxHostIdleConns:      1000,  // no limit
		ResponseHeaderTimeout: 5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		Timeout:               0,
	})

	DefaultAWSHttpClientLongTimeout = NewHTTPClientWithTimeouts(HTTPClientSettings{
		ConnectTimeout:        5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ConnKeepAlive:         30 * time.Second,
		MaxAllIdleConns:       10000, // no limit
		MaxHostIdleConns:      1000,  // no limit
		ResponseHeaderTimeout: 5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		Timeout:               1 * time.Minute,
	})
)
