package dokypay

import (
	"testing"
)

func TestDokypaySign(t *testing.T) {
	// 使用dokypay 文档上的测试示例
	sign := DokypaySign(map[string]interface{}{
		"amount":      "12.01",
		"appId":       "1000000126",
		"country":     "ID",
		"currency":    "USD",
		"description": "这是一个测试的商品",
		"merTransNo":  "mtn8888888888",
		"notifyUrl":   "http://yoursite.com/notifyurl",
		"pmId":        "doku",
		"prodName":    "southeast.asia",
		"returnUrl":   "http://yoursite.com/returnurl",
		"version":     "1.0",
	}, "5f190aa12f6442e0be4efca58680355b")
	if sign != "d11d877c0f435f2f8a263eca22559af47fb5679b7b235483f2b6767cd9d1ce22" {
		t.Error("not match")
	}
}
