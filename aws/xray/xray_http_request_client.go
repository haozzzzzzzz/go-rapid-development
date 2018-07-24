package xray

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/haozzzzzzzz/go-rapid-development/utils/http"
)

var RequestClientWithXray = xray.Client(http.RequestClient)
