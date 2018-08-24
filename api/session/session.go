package session

import (
	"time"

	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/sirupsen/logrus"
)

type AppSession struct {
	Ctx *ginbuilder.Context

	TraceId  string `json:"trace_id"`
	ClientIP string `json:"client_ip"`

	RequestData struct {
		AppVersion     string `json:"app_version"`
		AppVersionCode uint32 `json:"app_version_code"`
		DeviceId       string `json:"device_id"`
		AppType        string `json:"app_type"`
		ProductId      string `json:"product_id"`
	} `json:"request_data"`

	StartTime    time.Time
	EndTime      time.Time
	ExecDuration time.Duration
}

func (m *AppSession) GetLoggerFields() logrus.Fields {
	return logrus.Fields{
		"trace_id":         m.TraceId,
		"client_ip":        m.ClientIP,
		"app_version_code": m.RequestData.AppVersionCode,
		"app_device_id":    m.RequestData.DeviceId,
		"app_product_id":   m.RequestData.ProductId,
	}
}

type ManageSession struct {
	Ctx *ginbuilder.Context

	TraceId  string `json:"trace_id"`
	ClientIP string `json:"client_ip"`

	RequestData struct {
		UserId   string `json:"user_id" form:"user_id" binding:"required"`   // 用户ID
		UserName string `json:"username" form:"username" binding:"required"` // 用户名
	} `json:"request_data"`

	StartTime    time.Time
	EndTime      time.Time
	ExecDuration time.Duration
}

func (m *ManageSession) GetLoggerFields() logrus.Fields {
	return logrus.Fields{
		"trace_id":         m.TraceId,
		"client_ip":        m.ClientIP,
		"manage_user_id":   m.RequestData.UserId,
		"manage_user_name": m.RequestData.UserName,
	}
}
