package session

import (
	"time"

	"github.com/haozzzzzzzz/go-rapid-development/api/code"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

type AppSession struct {
	Ctx         *ginbuilder.Context
	RequestData struct {
		AppVersion     string `json:"app_version"`
		AppVersionCode uint32 `json:"app_version_code"`
		DeviceId       string `json:"device_id"`
		AppType        string `json:"app_type"`
		ProductId      string `json:"product_id"`
	} `json:"request_data"`

	ResponseData struct {
		ReturnCode *code.ApiCode `json:"return_code"`
	} `json:"response_data"`

	StartTime    time.Time
	EndTime      time.Time
	ExecDuration time.Duration
}

type ManageSession struct {
	Ctx         *ginbuilder.Context
	RequestData struct {
		UserId   string `json:"user_id" form:"user_id" binding:"required"`   // 用户ID
		UserName string `json:"username" form:"username" binding:"required"` // 用户名
	} `json:"request_data"`
	ResponseData struct {
		ReturnCode *code.ApiCode `json:"return_code"`
	} `json:"response_data"`

	StartTime    time.Time
	EndTime      time.Time
	ExecDuration time.Duration
}
