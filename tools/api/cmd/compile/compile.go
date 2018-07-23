package compile

import (
	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
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
			err = apiParser.ParseRouter()
			if nil != err {
				logrus.Errorf("parse router failed. %s.", err)
				return
			}

			// goimports
			goimports.DoGoImports([]string{serviceDir}, true)

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&serviceDir, "path", "p", "./", "service path")

	return cmd
}
