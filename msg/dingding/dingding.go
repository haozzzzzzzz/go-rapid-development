package dingding

import (
	"context"
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/v2/http"
	"github.com/haozzzzzzzz/go-rapid-development/v2/limiting/current_limiting"
	"github.com/haozzzzzzzz/go-rapid-development/v2/limiting/store"
	"github.com/haozzzzzzzz/go-rapid-development/v2/task"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/uerrors"
	"github.com/sirupsen/logrus"
	http2 "net/http"
	"reflect"
	"time"
)

// dingding alert with current limiting
type LimitingDingdingAlert struct {
	ApiUrl   string
	limiting *current_limiting.MinuteFrequencyLimiting
}

func NewLimitingDingdingAlert(apiUrl string) *LimitingDingdingAlert {
	return &LimitingDingdingAlert{
		ApiUrl: apiUrl,
	}
}

func (m *LimitingDingdingAlert) Start(
	workerNumber uint,
) (err error) {
	if m.limiting != nil {
		return
	}

	m.limiting = &current_limiting.MinuteFrequencyLimiting{
		WorkerNumber: workerNumber,
		Times:        20,
		WaitInterval: time.Second,
		MaxBatchSize: 10,
		Store:        store.NewMemoryStore(),
		Handler:      m.send,
	}

	err = m.limiting.Start()
	if nil != err {
		logrus.Errorf("start limiting failed. error: %s.", err)
		return
	}
	return
}

type DingdingMsg struct {
	Service    string
	Content    string
	CreateTime time.Time
}

func (m *LimitingDingdingAlert) Send(msg *DingdingMsg) (err error) {
	err = m.limiting.AcceptData(msg)
	if nil != err {
		logrus.Errorf("limiting accept data failed. error: %s.", err)
		return
	}
	return
}

func (m *LimitingDingdingAlert) send(datas ...interface{}) (err error) {
	if len(datas) == 0 {
		return
	}

	var contents string
	for _, data := range datas {
		msg, ok := data.(*DingdingMsg)
		if !ok {
			err = uerrors.Newf("convert to *DingdingMsg failed. type: %s", reflect.TypeOf(data))
			return
		}

		content := fmt.Sprintf("[%s]\n%s\n%s", msg.Service, msg.Content, msg.CreateTime)
		if contents != "" {
			contents += "\n" + content
		} else {
			contents = content
		}
	}

	logrus.Infof("dingding send : %#v", contents)
	ctx, cancel := task.NewBackgroundContext("dingding_send")
	defer func() {
		cancel(err)
	}()

	err = SendDingdingMsg(m.ApiUrl, ctx, http.LongTimeoutRequestClient, contents)
	if err != nil {
		logrus.Errorf("send dingding msg failed. error: %s", err)
		return
	}

	return
}

func SendDingdingMsg(
	apiUrl string,
	ctx context.Context,
	httpClient *http2.Client,
	contents string,
) (err error) {
	req, err := http.NewRequest(apiUrl, ctx, httpClient)
	if nil != err {
		logrus.Errorf("new http request failed. error: %s.", err)
		return
	}

	resp := &struct {
		Errmsg  string `json:"errmsg"`
		Errcode uint32 `json:"errcode"`
	}{}
	err = req.PostJson(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": contents,
		},
		"at": map[string]interface{}{
			"isAtAll": false,
		},
	}, resp)
	if nil != err {
		logrus.Errorf("http request post json failed. error: %s.", err)
		return
	}

	if resp.Errcode != 0 {
		err = uerrors.Newf("dingding response error. errcode: %d, errmsg: %s", resp.Errcode, resp.Errmsg)
	}
	return
}
