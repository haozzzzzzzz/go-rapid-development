package main

import (
	_ "github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/app/first_init"
	// 空一行

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/app/api"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/app/constant"

	_ "service/common/session/app"

	_ "github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/common/config"

	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"service/common/operations"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/common/config"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	flags.StringVarP(&runParams.Host, "ip", "i", "", "serve host ip")
	flags.StringVarP(&runParams.Port, "port", "p", "18200", "serve port")

	if err := mainCmd.Execute(); err != nil {
		logrus.Println(err)
		os.Exit(1)
	}

}

type RunParams struct {
	Host string
	Port string
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

	// 开启跨域
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))

	// bind xray
	engine.Use(xray.XRayGinMiddleware(serviceName))

	// bind prometheus
	metricsPath := "/api/test_doc/metrics"
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
