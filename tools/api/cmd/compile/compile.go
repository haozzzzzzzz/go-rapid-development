package compile

import (
	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandApiCompile() *cobra.Command {
	var projectPath string
	var cmd = &cobra.Command{
		Use:   "compile",
		Short: "api project compilation",
		Run: func(cmd *cobra.Command, args []string) {
			if projectPath == "" {
				logrus.Errorf("need path")
				return
			}

			var err error
			projectPath, err = filepath.Abs(projectPath)
			if nil != err {
				logrus.Errorf("get absolute project path failed. \ns%s.", err)
				return
			}

			err = api.NewApiParser(projectPath).MapApi()
			if nil != err {
				logrus.Errorf("mapping api failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectPath, "path", "p", "./", "project path")

	return cmd
}
