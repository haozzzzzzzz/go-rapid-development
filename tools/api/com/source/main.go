package source

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateMain(params *GenerateParams) (err error) {
	projDir := m.ProjectDir

	// generate main file
	mainFilePath := fmt.Sprintf("%s/main.go", projDir)
	newMainFileText := strings.Replace(mainFileText, "$HOST$", params.Host, -1)
	newMainFileText = strings.Replace(newMainFileText, "$PORT$", params.Port, -1)
	err = ioutil.WriteFile(mainFilePath, []byte(newMainFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("new project main file failed. %s.", err)
		return
	}

	return
}

// main.go
var mainFileText = `package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	// metric
	defer func() {
		if err := recover(); err != nil {
			metrics.SERVICE_TIMES_COUNTER_PANIC.Inc()
		}
	}()

	runParams := &RunParams{}
	mainCmd := &cobra.Command{
		Long: fmt.Sprintf("%s service", constant.ServiceName),
		Run: func(cmd *cobra.Command, args []string) {
			Run(runParams)
		},
	}

	flags := mainCmd.Flags()
	flags.StringVarP(&runParams.Host, "ip", "i", "$HOST$", "serve host ip")
	flags.StringVarP(&runParams.Port, "port", "p", "$PORT$", "serve port")
	flags.StringVarP(&runParams.Stage, "stage", "s", "test", "deploy stage. dev、test、pre、prod")

	if err := mainCmd.Execute(); err != nil {
		logrus.Println(err)
		os.Exit(1)
	}

}

type RunParams struct {
	Host  string
	Port  string
	Stage string
}

func Run(runParams *RunParams) {

	engine := ginbuilder.GetEngine()

	// bind xray
	engine.Use(xray.XRayGinMiddleware(fmt.Sprintf("%s_%s", runParams.Stage, constant.ServiceName)))

	// bind prometheus
	engine.GET(fmt.Sprintf("/%s/metrics", constant.ServiceName), func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Println(err)
			}
		}()

		promhttp.Handler().ServeHTTP(context.Writer, context.Request)

	})

	api.BindRouters(engine)
	engine.Run(fmt.Sprintf("%s:%s", runParams.Host, runParams.Port))

}
`
