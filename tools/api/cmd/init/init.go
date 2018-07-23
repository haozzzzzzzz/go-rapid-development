package init

import (
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	service2 "github.com/haozzzzzzzz/go-rapid-development/tools/api/com/service"
)

// 初始化服务框架
func CommandApiInitFramework() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "init_framework",
		Short: "api project framework initialization",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return
}

// 初始化API服务
func CommandApiInitService() *cobra.Command {
	var config project.ServiceConfigFormat
	var params service2.GenerateParams
	var cmd = &cobra.Command{
		Use:   "init_service",
		Short: "api service initialization",
		Run: func(cmd *cobra.Command, args []string) {
			err := validator.New().Struct(&config)
			if nil != err {
				logrus.Errorf("validate init config failed. %s.", err)
				return
			}

			absServiceDir, err := filepath.Abs(config.ServiceDir)
			if nil != err {
				logrus.Errorf("get absolute file path %q failed. %s.", absServiceDir, err)
				return
			}
			config.ServiceDir = absServiceDir

			service := &project.Service{
				Config: &config,
			}
			err = service.Init()
			if nil != err {
				logrus.Errorf("save service config failed. %s.", err)
				return
			}

			// init api service files
			apiServiceSource := service2.NewApiServiceSource(service)
			err = apiServiceSource.Generate(&params)
			if nil != err {
				logrus.Errorf("generate api service source failed. %s.", err)
				return
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&config.Name, "name", "n", "", "api service name")
	flags.StringVarP(&config.ServiceDir, "path", "p", "./", "api service directory path")
	flags.StringVarP(&config.Description, "description", "d", "api", "api service description")
	flags.StringVarP(&params.Host, "host", "H", "", "api serve host")
	flags.StringVarP(&params.Port, "port", "P", "18100", "api serve port")

	return cmd
}
