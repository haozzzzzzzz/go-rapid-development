package dokypay

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/fatih/structs"
	"github.com/haozzzzzzzz/go-rapid-development/api/client"
	"github.com/haozzzzzzzz/go-rapid-development/utils/str"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
)

/**
https://www.dokypay.com/api.html#153042ed9e
*/

type Dokypay struct {
	Ctx         context.Context
	HttpClient  *http.Client
	MerchantId  string // 商家ID
	MerchantKey string // 商家密钥
	AppId       string // 应用ID
	AppKey      string // 应用密钥
	ApiPrefix   string
}

func NewDokypay(
	merchantId string,
	merchantKey string,
	appId string,
	appKey string,
	apiPrefix string,
) *Dokypay {
	return &Dokypay{
		MerchantId:  merchantId,
		MerchantKey: merchantKey,
		AppId:       appId,
		AppKey:      appKey,
		ApiPrefix:   apiPrefix,
	}
}

// api
const ApiVersion string = "1.0"

type ResponseData struct {
	Code     string      `json:"code"`
	Msg      string      `json:"msg"`
	Timetamp string      `json:"timetamp"`
	Data     interface{} `json:"data"`
}

func (m *ResponseData) CheckCode() (err error) {
	if m.Code == "200" {
		return
	}

	err = uerrors.Newf("dokyapi response error. code: %s, msg: %s, data: %#v", m.Code, m.Msg, m)
	return
}

const TradeStatusPending string = "pending"
const TradeStatusProcessing string = "processing"
const TradeStatusSuccess string = "success"
const TradeStatusFailure string = "failure"
const TradeStatusTimeout string = "timeout"
const TradeStatusCancel string = "cancel"

type PayoutAsyncNotifyPostData struct {
	ResultCode  string `json:"result_code" form:"result_code" structs:"resultCode" mapstructure:"resultCode"`
	Message     string `json:"message" form:"message" structs:"message" mapstructure:"message"`
	TradeNo     string `json:"tradeNo" form:"tradeNo" structs:"tradeNo" mapstructure:"tradeNo" binding:"required"`
	MerTransNo  string `json:"merTransNo" form:"merTransNo" structs:"merTransNo" mapstructure:"merTransNo" binding:"required"`
	Currency    string `json:"currency" form:"currency" structs:"currency" mapstructure:"currency" binding:"required"`
	Amount      string `json:"amount" form:"amount" structs:"amount" mapstructure:"amount" binding:"required"`
	TotalFee    string `json:"totalFee" form:"totalFee" structs:"totalFee" mapstructure:"totalFee" binding:"required"`
	TradeStatus string `json:"tradeStatus" form:"tradeStatus" structs:"tradeStatus" mapstructure:"tradeStatus" binding:"required"`
	CreateTime  string `json:"createTime" form:"createTime" structs:"createTime" mapstructure:"createTime" binding:"required"` // 北京时间 GMT +8
	UpdateTime  string `json:"updateTime" form:"updateTime" structs:"updateTime" mapstructure:"updateTime" binding:"required"` // 北京时间
	Sign        string `json:"sign" form:"sign" structs:"sign" mapstructure:"sign" binding:"required"`
}

type PayoutRequestData struct {
	Amount      string      `json:"amount" form:"amount" structs:"amount" binding:"required"`
	AppId       string      `json:"appId" form:"appId" structs:"appId" binding:"required"`
	Currency    string      `json:"currency" form:"currency" structs:"currency" binding:"required"`
	Description string      `json:"description" form:"description" structs:"description"`
	MerTransNo  string      `json:"merTransNo" form:"merTransNo" structs:"merTransNo" binding:"required"`
	NotifyUrl   string      `json:"notifyUrl" form:"notifyUrl" structs:"notifyUrl" binding:"required"`
	PmId        string      `json:"pmId" form:"pmId" structs:"pmId" binding:"required"`
	ProdName    string      `json:"prodName" form:"prodName" structs:"prodName" binding:"required"`
	Version     string      `json:"version" form:"version" structs:"version" binding:"required"`
	Sign        string      `json:"sign" form:"sign" structs:"sign" binding:"required"`
	ExtInfo     interface{} `json:"extInfo" form:"extInfo" structs:"extInfo" binding:"required"`
}

const ResultCodeSuccess = "0000"

type PayoutResponseResult struct {
	ResultCode string `json:"resultCode" structs:"resultCode"`
	Message    string `json:"message" structs:"message"`
	Amount     string `json:"amount" structs:"amount"`
	Currency   string `json:"currency" structs:"currency"`
	MerTransNo string `json:"merTransNo" structs:"merTransNo"`
	TradeNo    string `json:"tradeNo" structs:"tradeNo"`
	Sign       string `json:"sign" structs:"sign"`
}

const PmIdPaytmWalletPayout string = "paytm.wallet.payout"

type PaytmExtInfo struct {
	PayeeEmail  string `json:"payeeEmail,omitempty" structs:"payee_email,omitempty"`
	PayeeMobile string `json:"payeeMobile,omitempty" structs:"payee_mobile,omitempty"`
}

func (m *PaytmExtInfo) Map() map[string]interface{} {
	return structs.Map(m)
}

const ProdNameSouthEastAsiaPayout string = "southeast.asia.payout"

/**
{
    "amount": "12.01",
    "appId": "1000000126",
    "country": "ID",
    "currency": "USD",
    "description": "这是一个测试的商品",
    "merTransNo": "mtn8888888888",
    "notifyUrl": "http://yoursite.com/notifyurl",
    "pmId": "doku",
    "prodName": "southeast.asia",
    "returnUrl": "http://yoursite.com/returnurl",
    "version": "1.0"
}
商户的Key如下： 5f190aa12f6442e0be4efca58680355b
则签名前拼接字符串为： amount=12.01&appId=1000000126&country=ID&currency=USD&description=这是一个测试的商品&merTransNo=mtn8888888888&notifyUrl=http://yoursite.com/notifyurl&pmId=doku&prodName=southeast.asia&returnUrl=http://yoursite.com/returnurl&version=1.0&key=5f190aa12f6442e0be4efca58680355b
获取的签名值为： d11d877c0f435f2f8a263eca22559af47fb5679b7b235483f2b6767cd9d1ce22
*/

// 只对string类型的第一级参数进行签名
func DokypaySign(params map[string]interface{}, signKey string) (strSign string) {
	delete(params, "sign")

	keys := make([]string, 0)
	for key, param := range params {
		_, ok := param.(string)
		if !ok {
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	pairs := make([]string, 0)
	for _, key := range keys {
		strParam := params[key].(string)
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, strParam))
	}

	pairs = append(pairs, fmt.Sprintf("key=%s", signKey))
	originStr := strings.Join(pairs, "&")
	strSign = str.Sha256([]byte(originStr))
	return
}

func (m *Dokypay) PayoutApi(
	ctx context.Context,
	httpClient *http.Client,
	data *PayoutRequestData,
) (
	result *PayoutResponseResult,
	err error,
) {
	result = &PayoutResponseResult{}
	resp := &ResponseData{
		Data: result,
	}

	data.Sign = m.PayoutRequestDataSign(data)

	uri := "/payout"
	reqClient := client.NewClient(ctx, httpClient, m.ApiPrefix)
	err = reqClient.Post(uri, resp, nil, nil, data)
	if nil != err {
		logrus.Errorf("post dokypay payout api failed. error: %s.", err)
		return
	}

	err = resp.CheckCode()
	if nil != err {
		logrus.Errorf("check response code failed. uri: %s, error: %s.", uri, err)
		return
	}

	if result.ResultCode != "0000" { // 下单失败
		err = uerrors.Newf("payout api add order failed. return_code: %s, msg: %s", result.ResultCode, result.Message)
		return
	}

	// check sign
	calcSign := m.PayoutResponseResultSign(result)
	if calcSign != result.Sign {
		err = uerrors.Newf("payout return data sign not match. calc_sign: %s, ret_sign: %s, result: %#v", calcSign, result.Sign, result)
		return
	}

	return
}

// 不同请求的签名key可能是AppKey或者MerchantKey
func (m *Dokypay) PayoutRequestDataSign(data *PayoutRequestData) string {
	return DokypaySign(structs.Map(data), m.AppKey)
}

func (m *Dokypay) PayoutResponseResultSign(data *PayoutResponseResult) string {
	return DokypaySign(structs.Map(data), m.AppKey)
}

func (m *Dokypay) PayoutAsyncNotifyPostDataSign(data map[string]interface{}) string {
	return DokypaySign(data, m.AppKey)
}
