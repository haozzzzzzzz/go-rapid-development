package service

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateMain(params *GenerateParams) (err error) {
	serviceDir := m.ServiceDir

	// generate main file
	mainFilePath := fmt.Sprintf("%s/main.go", serviceDir)

	var metricPanicCounter string
	serviceType := project.ServiceType(m.Service.Config.Type)

	switch serviceType {
	case project.ServiceTypeApp:
		metricPanicCounter = "SERVICE_TIMES_COUNTER_APP_PANIC"
	case project.ServiceTypeManage:
		metricPanicCounter = "SERVICE_TIMES_COUNTER_MANAGE_PANIC"
	case project.ServiceTypeRPC:
		metricPanicCounter = "SERVICE_TIMES_COUNTER_RPC_PANIC"
	default:
		err = uerrors.Newf("unknown service type. %s", serviceType)
		return
	}

	newMainFileText := strings.Replace(mainFileText, "$METRIC_PANIC_COUNTER$", metricPanicCounter, -1)
	newMainFileText = strings.Replace(newMainFileText, "$HOST$", params.Host, -1)
	newMainFileText = strings.Replace(newMainFileText, "$PORT$", params.Port, -1)
	err = ioutil.WriteFile(mainFilePath, []byte(newMainFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("new service main file failed. %s.", err)
		return
	}

	return
}

// main.go
var mainFileText = `package main

import (
	// TODO 这里import first_init
	// 空一行

	// TODO 这里需要import config
	// TODO 这里需要导入session
	
	"fmt"
	"math/rand"
	"os"
	"time"
	"log"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"service/common/operations"
)

func main() {
	runParams := &RunParams{}
	mainCmd := &cobra.Command{
		Long: fmt.Sprintf("%s service", constant.ServiceName),
		Run: func(cmd *cobra.Command, args []string) {
			err := Run(runParams)
			if nil != err {
				logrus.Errorf("run service failed. %s.", err)
				return
			}
		},
	}

	flags := mainCmd.Flags()
	flags.StringVarP(&runParams.Host, "ip", "i", "$HOST$", "serve host ip")
	flags.StringVarP(&runParams.Port, "port", "p", "$PORT$", "serve port")

	if err := mainCmd.Execute(); err != nil {
		logrus.Println(err)
		os.Exit(1)
	}

}

type RunParams struct {
	Host  string
	Port  string
}

func Run(runParams *RunParams) (err error) {
	serviceName := config.EnvConfig.WithStagePrefix(constant.ServiceName)
	rand.Seed(time.Now().Unix())

	if config.EnvConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := ginbuilder.DefaultEngine()

	// bind xray
	engine.Use(xray.XRayGinMiddleware(serviceName))

	// bind prometheus
	// TODO 需要添加前缀
	metricsPath := "/api/__/metrics"
	engine.GET(metricsPath, func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Println(err)
			}
		}()

		promhttp.Handler().ServeHTTP(context.Writer, context.Request)

	})

	api.BindRouters(engine)

	log.Printf("Running %s on %s:%s\n", serviceName, runParams.Host, runParams.Port)
	address := fmt.Sprintf("%s:%s", runParams.Host, runParams.Port)

	// 注册
	opTool := &operations.OperationsTool{
		ServiceName:      constant.ServiceName,
		MetricsPath:      metricsPath,
		PrometheusTarget: fmt.Sprintf("%s:%s", config.AWSEc2InstanceIdentifyDocument.PrivateIP, runParams.Port),
	}
	err = opTool.RegisterPrometheus()
	if nil != err {
		logrus.Errorf("register prometheus failed. %s", err)
		err = nil // 容错
	}

	defer func() {
		err = opTool.UnregisterPrometheus()
		if nil != err {
			logrus.Errorf("unregister prometheus failed. %s", err)
			return
		}
	}()

	endless.DefaultReadTimeOut = 10 * time.Second
	endless.DefaultWriteTimeOut = 10 * time.Second
	err = endless.ListenAndServe(address, engine)
	if nil != err {
		logrus.Errorf("start listening and serving http on %s failed. %s", address, err)
		return
	}

	return
}
`
