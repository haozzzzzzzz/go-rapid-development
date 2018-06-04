package init

import (
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/common/proj"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// 初始化API项目
func CommandApiInit() *cobra.Command {
	var config proj.ProjectConfigFormat
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "api project initialization",
		Run: func(cmd *cobra.Command, args []string) {
			err := validator.New().Struct(&config)
			if nil != err {
				logrus.Errorf("validate init config failed. %s.", err)
				return
			}

			absProjectDir, err := filepath.Abs(config.ProjectDir)
			if nil != err {
				logrus.Errorf("get absolute file path %q failed. %s.", absProjectDir, err)
				return
			}
			config.ProjectDir = absProjectDir

			project := &proj.Project{
				Config: &config,
			}
			err = project.Save()
			if nil != err {
				logrus.Errorf("save project config failed. %s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&config.Name, "name", "n", "", "api project name")
	flags.StringVarP(&config.ProjectDir, "path", "p", "./", "api project directory path")
	flags.StringVarP(&config.Description, "description", "d", "api", "api project description")

	return cmd
}
