package compile

import (
	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/parser"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandApiCompile() *cobra.Command {
	var serviceDir string
	var cmd = &cobra.Command{
		Use:   "compile",
		Short: "api service compilation",
		Run: func(cmd *cobra.Command, args []string) {
			if serviceDir == "" {
				logrus.Errorf("service dir required")
				return
			}

			var err error
			serviceDir, err = filepath.Abs(serviceDir)
			if nil != err {
				logrus.Errorf("get absolute service path failed. \ns%s.", err)
				return
			}

			// service
			service, err := project.LoadService(serviceDir)
			if nil != err {
				logrus.Errorf("load service failed. %s.", err)
				return
			}

			// api parser
			apiParser := parser.NewApiParser(service)
			apis, err := apiParser.ScanApis(false, false)
			if nil != err {
				logrus.Errorf("Scan api failed. \n%s.", err)
				return
			}

			err = apiParser.GenerateRoutersSourceFile(apis)
			if nil != err {
				logrus.Errorf("generate routers source file failed. %s.", err)
				return
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&serviceDir, "path", "p", "./", "service path")

	return cmd
}
