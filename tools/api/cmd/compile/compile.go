package compile

import (
	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/parser"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandApiCompile() *cobra.Command {
	var projectDir string
	var cmd = &cobra.Command{
		Use:   "compile",
		Short: "api project compilation",
		Run: func(cmd *cobra.Command, args []string) {
			if projectDir == "" {
				logrus.Errorf("project dir required")
				return
			}

			var err error
			projectDir, err = filepath.Abs(projectDir)
			if nil != err {
				logrus.Errorf("get absolute project path failed. \ns%s.", err)
				return
			}

			// project
			project, err := proj.LoadProject(projectDir)
			if nil != err {
				logrus.Errorf("load project failed. %s.", err)
				return
			}

			// api parser
			apiParser := parser.NewApiParser(project)
			err = apiParser.ParseRouter()
			if nil != err {
				logrus.Errorf("parse router failed. %s.", err)
				return
			}

			// do go imports
			goimports.DoGoImports([]string{projectDir}, true)

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectDir, "path", "p", "./", "project path")

	return cmd
}
