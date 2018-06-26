package client

import (
	"context"
	"fmt"
	"service/VideoBuddyConfigApi/common/db/model"
	"testing"
)

func TestClient_Get(t *testing.T) {
	client := Client{
		ctx:        context.Background(),
		HttpClient: nil,
	}

	urlPrefix := "http://127.0.0.1:18105"
	urlPath := "/config_manage/v1/condition_group/attachment/filter_list/:config_type"
	iRespData := &struct {
		ReturnCode uint32      `json:"ret"`
		Message    string      `json:"msg"`
		Data       interface{} `json:"data"`
	}{}

	iPathData := &struct {
		ConfigType string `json:"config_type" form:"config_type"`
	}{
		ConfigType: "home",
	}

	iQueryData := &struct {
		Offset          uint32                          `json:"offset" form:"offset"` // attachment id
		JumpPage        uint32                          `json:"jump_page" form:"jump_page"`
		Limit           uint8                           `json:"limit" form:"limit" binding:"required"`
		OnlineState     model.AttachmentOnlineStateType `json:"online_state" form:"online_state"`                         // 上线状态 0：未上线；1：已上线
		FilterEffective uint8                           `json:"filter_effective" form:"filter_effective" binding:"lte=2"` // 生效状态过滤。0：不过滤；1：已生效；2：未生效
		UserId          string                          `json:"user_id" form:"user_id"`
		Username        string                          `json:"username" form:"username"`
	}{
		Offset:          0,
		JumpPage:        0,
		Limit:           10,
		OnlineState:     model.AttachmentOnlineStateOnline,
		FilterEffective: 1,
		UserId:          "1",
		Username:        "luohao",
	}

	err := client.Get(urlPrefix, urlPath, iRespData, iPathData, iQueryData)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(iRespData)
}