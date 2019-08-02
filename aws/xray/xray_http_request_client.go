package xray

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/haozzzzzzzz/go-rapid-development/http"
)

var RequestClientWithXray = xray.Client(http.ShortTimeoutRequestClient)
var LongTimeoutRequestClientWithXray = xray.Client(http.LongTimeoutRequestClient)
var NoTimeoutRequestClientWithXray = xray.Client(http.NoTimeoutRequestClient)
