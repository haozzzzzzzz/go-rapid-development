package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateApi() (err error) {
	projDir := m.ProjectDir

	// generate api folder
	apiDir := fmt.Sprintf("%s/api", projDir)
	err = os.MkdirAll(apiDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project api dir %q failed. %s.", apiDir, err)
		return
	}

	// routers
	routersFilePath := fmt.Sprintf("%s/routers.go", apiDir)
	err = ioutil.WriteFile(routersFilePath, []byte(routersFileText), proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("write api/routers.go failed. \n%s.", err)
		return
	}

	// session
	sessionDir := fmt.Sprintf("%s/session", apiDir)
	err = os.MkdirAll(sessionDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make api session dir %q failed. %s.", err)
		return
	}

	sessionFilePath := fmt.Sprintf("%s/session.go", sessionDir)
	err = ioutil.WriteFile(sessionFilePath, []byte(sessionFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write session file %q failed. %s.", sessionFilePath, err)
		return
	}

	return
}

// routers.go
var routersFileText = `package api

import (
	"github.com/gin-gonic/gin"
)

// 注意：BindRouters函数体内不能自定义添加任何声明，由lambda-build compile api命令生成api绑定声明
func BindRouters(engine *gin.Engine) (err error) {
	return
}
`

var sessionFileText = `package session

import (
	"time"

	"fmt"

	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/api/code"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

type ApiSession struct {
	RequestData struct {
		AppVersion     string ` + "`json:\"app_version\"`" + `
		AppVersionCode string ` + "`json:\"app_version_code\"`" + `
		DeviceId       string ` + "`json:\"device_id\"`" + `
		AppType        string ` + "`json:\"app_type\"`" + `
		ProductId      string ` + "`json:\"product_id\"`" + `
	} ` + "`json:\"request_data\"`" + `

	ResponseData struct {
		ReturnCode *ginbuilder.ReturnCode ` + "`json:\"return_code\"`" + `
	} ` + "`json:\"response_data\"`" + `

	StartTime    time.Time
	EndTime      time.Time
	ExecDuration time.Duration
}

func (m *ApiSession) BeforeHandle(ctx *ginbuilder.Context, method string, uri string) {
	m.StartTime = time.Now()

	// metrics
	metrics.API_EXEC_TIMES_COUNTER_TOTAL.Inc()
	metrics.API_URI_CALL_TIMES_COUNTER_VEC.WithLabelValues(uri).Inc()
}

func (m *ApiSession) AfterHandle(err error) {
	m.EndTime = time.Now()
	m.ExecDuration = m.EndTime.Sub(m.StartTime)

	// metrics
	metrics.API_SPENT_TIME_SUMMARY.Observe(m.ExecDuration.Seconds() * 1000)
	if m.ResponseData.ReturnCode != nil {
		retCode := m.ResponseData.ReturnCode
		metrics.API_RETURN_CODE_COUNTER_VEC.WithLabelValues(fmt.Sprintf("%d", retCode.Code)).Inc()

		if retCode.Code != code.CodeSuccess.Code {
			metrics.API_EXEC_TIMES_COUNTER_ABNORMAL.Inc()
		}
	}
}

func (m *ApiSession) Panic(err interface{}) {
	metrics.API_EXEC_TIMES_COUNTER_PANIC.Inc()
}

func (m *ApiSession) SetReturnCode(code *ginbuilder.ReturnCode) {
	m.ResponseData.ReturnCode = code
}

func GetSession(ctx *ginbuilder.Context) (session *ApiSession) {
	session = ctx.Session.(*ApiSession)
	return
}

func SessionBuilder(ctx *ginbuilder.Context) (err error) {
	ses := &ApiSession{}
	ctx.Session = ses

	reqHeader := ctx.GinContext.Request.Header
	ses.RequestData.AppVersion = reqHeader.Get("Version-Name")
	ses.RequestData.AppVersionCode = strings.Trim(reqHeader.Get("Version-Code"), "")
	ses.RequestData.DeviceId = reqHeader.Get("Device-Id")
	ses.RequestData.AppType = reqHeader.Get("App-Type")
	ses.RequestData.ProductId = reqHeader.Get("Product-Id")

	return
}

func init() {
	ginbuilder.BindSessionBuilder(SessionBuilder)
}
`
