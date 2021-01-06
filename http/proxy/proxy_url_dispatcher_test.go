package proxy

import (
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestProxyUrls_ReplaceUrls(t *testing.T) {
	proxyUrls := NewProxyUrls()
	for i := 1; i < 10; i++ {
		proxyUrls.AddUrl(NewProxyUrl(fmt.Sprintf("http://%d.com", i), rate.NewLimiter(rate.Limit(2), i)))
	}

	for i := 0; i < 100; i++ {
		success, strUrl, err := proxyUrls.SelectUrl(1 * time.Second)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(success, strUrl)
	}
}
