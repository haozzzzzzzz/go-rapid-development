package service

import (
	"fmt"
	"os"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
)

func (m *ServiceSource) generateApi() (err error) {
	serviceDir := m.ServiceDir

	// generate api folder
	apiDir := fmt.Sprintf("%s/api", serviceDir)
	err = os.MkdirAll(apiDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service api dir %q failed. %s.", apiDir, err)
		return
	}

	// routers
	routersFilePath := fmt.Sprintf("%s/routers.go", apiDir)
	err = ioutil.WriteFile(routersFilePath, []byte(routersFileText), project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("write api/routers.go failed. \n%s.", err)
		return
	}

	// session
	sessionDir := fmt.Sprintf("%s/session", apiDir)
	err = os.MkdirAll(sessionDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make api session dir %q failed. %s.", err)
		return
	}

	sessionFilePath := fmt.Sprintf("%s/session.go", sessionDir)

	var sessionFileText string
	serviceType := project.ServiceType(m.Service.Config.Type)
	switch serviceType {
	case project.ServiceTypeApp:
		sessionFileText = appSessionFileText
	case project.ServiceTypeManage:
		sessionFileText = manageSessionFileText
	default:
		err =uerrors.Newf("unknown service type %s", serviceType)
		return
	}

	err = ioutil.WriteFile(sessionFilePath, []byte(sessionFileText), project.ProjectFileMode)
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

var appSessionFileText = `package session

import (
	"strconv"
	"time"

	"fmt"


	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/api/code"
	"github.com/haozzzzzzzz/go-rapid-development/api/session"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/sirupsen/logrus"
)

type AppSession struct {
	session.AppSession
}

func (m *AppSession) BeforeHandle(ctx *ginbuilder.Context, method string, uri string) {
	m.StartTime = time.Now()

	// metrics
	metrics.API_EXEC_TIMES_COUNTER_APP_TOTAL.Inc()
	metrics.API_URI_CALL_TIMES_COUNTER_VEC.WithLabelValues(metrics.Ec2InstanceId, "app", uri).Inc()
}

func (m *AppSession) AfterHandle(err error) {
	m.EndTime = time.Now()
	m.ExecDuration = m.EndTime.Sub(m.StartTime)

	// metrics
	metrics.API_SPENT_TIME_SUMMARY_APP.Observe(m.ExecDuration.Seconds() * 1000)

	if m.ResponseData.ReturnCode != nil {
		retCode := m.ResponseData.ReturnCode
		metrics.API_RETURN_CODE_COUNTER_VEC.WithLabelValues(metrics.Ec2InstanceId, "app", fmt.Sprintf("%d", retCode.Code)).Inc()

		if retCode.Code != code.CodeSuccess.Code {
			metrics.API_EXEC_TIMES_COUNTER_APP_ABNORMAL.Inc()
		}
	}

}

func (m *AppSession) Panic(err interface{}) {
	metrics.API_EXEC_TIMES_COUNTER_APP_PANIC.Inc()
}

func (m *AppSession) SetReturnCode(code *code.ApiCode) {
	m.ResponseData.ReturnCode = code
}

func GetSession(ctx *ginbuilder.Context) (session *AppSession) {
	session = ctx.Session.(*AppSession)
	return
}

func SessionBuilder(ctx *ginbuilder.Context) (err error) {
	ses := &AppSession{}
	ctx.Session = ses

	reqHeader := ctx.GinContext.Request.Header

	ses.RequestData.AppVersion = reqHeader.Get("Version-Name")
	strVersionCode := strings.Trim(reqHeader.Get("Version-Code"), "")

	if strVersionCode != "" {
		appVersionCode, errParse := strconv.ParseInt(strVersionCode, 10, 32)
		err = errParse
		if err != nil {
			logrus.Errorf("parse string version code to int failed. %s.", err)
			return
		}
		ses.RequestData.AppVersionCode = uint32(appVersionCode)

	}

	ses.RequestData.DeviceId = reqHeader.Get("Device-Id")
	ses.RequestData.AppType = reqHeader.Get("App-Type")
	ses.RequestData.ProductId = reqHeader.Get("Product-Id")

	return
}

func init() {
	ginbuilder.BindSessionBuilder(SessionBuilder)
}
`

var manageSessionFileText = `package session

import (
	"time"

	"fmt"

	"github.com/haozzzzzzzz/go-rapid-development/api/code"
	"github.com/haozzzzzzzz/go-rapid-development/api/session"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

type ManageSession struct {
	session.ManageSession
}

func (m *ManageSession) BeforeHandle(ctx *ginbuilder.Context, method string, uri string) {
	m.StartTime = time.Now()

	// metrics
	metrics.API_EXEC_TIMES_COUNTER_MANAGE_TOTAL.Inc()
	metrics.API_URI_CALL_TIMES_COUNTER_VEC.WithLabelValues(metrics.Ec2InstanceId, "manage", uri).Inc()

}

func (m *ManageSession) AfterHandle(err error) {
	m.EndTime = time.Now()
	m.ExecDuration = m.EndTime.Sub(m.StartTime)

	// metrics
	metrics.API_SPENT_TIME_SUMMARY_MANAGE.Observe(m.ExecDuration.Seconds() * 1000)
	if m.ResponseData.ReturnCode != nil {
		retCode := m.ResponseData.ReturnCode
		metrics.API_RETURN_CODE_COUNTER_VEC.WithLabelValues(metrics.Ec2InstanceId, "manage", fmt.Sprintf("%d", retCode.Code)).Inc()

		if retCode.Code != code.CodeSuccess.Code {
			metrics.API_EXEC_TIMES_COUNTER_MANAGE_ABNORMAL.Inc()
		}
	}
}

func (m *ManageSession) Panic(err interface{}) {
	metrics.API_EXEC_TIMES_COUNTER_MANAGE_PANIC.Inc()
}

func (m *ManageSession) SetReturnCode(code *code.ApiCode) {
	m.ResponseData.ReturnCode = code
}

func GetSession(ctx *ginbuilder.Context) (session *ManageSession) {
	session = ctx.Session.(*ManageSession)
	return
}

func SessionBuilder(ctx *ginbuilder.Context) (err error) {
	sess := &ManageSession{}
	ctx.Session = sess
	retCode, err := ctx.BindQueryData(&sess.RequestData)
	if nil != err {
		ctx.Errorf(retCode, "bind session query data failed. %s", err)
		return
	}

	return
}

func init() {
	ginbuilder.BindSessionBuilder(SessionBuilder)
}
`