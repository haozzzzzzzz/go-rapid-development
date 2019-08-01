package compile

import (
	"fmt"
	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/parser"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GenerateApiDoc() *cobra.Command {
	var serviceDir string
	var host string
	var version string
	var contactName string
	var noMod bool
	var cmd = &cobra.Command{
		Use:   "doc",
		Short: "api doc generate",
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
			apis, err := apiParser.ScanApis(true, !noMod)
			if nil != err {
				logrus.Errorf("Scan api failed. \n%s.", err)
				return
			}

			swaggerSpec := parser.NewSwaggerSpec()
			swaggerSpec.Info(
				service.Config.Name,
				service.Config.Description,
				version,
				contactName,
			)
			swaggerSpec.Host(host)
			swaggerSpec.Apis(apis)
			swaggerSpec.Schemes([]string{"http", "https"})
			err = swaggerSpec.ParseApis()
			err = swaggerSpec.SaveToFile(fmt.Sprintf("%s/.service/swagger.json", service.Config.ServiceDir))
			if nil != err {
				logrus.Errorf("save swagger spec to file failed. error: %s.", err)
				return
			}

			err = apiParser.SaveApisToFile(apis)
			if nil != err {
				logrus.Errorf("save apis to file failed. %s.", err)
				return
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&serviceDir, "path", "p", "./", "service path")
	flags.StringVarP(&host, "host", "H", "", "api host")
	flags.StringVarP(&version, "version", "v", "1.0", "api version")
	flags.StringVarP(&contactName, "contact_name", "c", "", "contact name")
	flags.BoolVarP(&noMod, "no_mod", "n", false, "use mod")
	return cmd
}

func GenerateCommentDoc() (cmd *cobra.Command) {
	var dir string
	var outputPath string
	var host string
	cmd = &cobra.Command{
		Use:   "com_doc",
		Short: "generate open api doc from comment",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if dir == "" || outputPath == "" {
				logrus.Errorf("invalid params")
				return
			}

			dir, err = filepath.Abs(dir)
			if nil != err {
				logrus.Errorf("get abs dir failed. dir: %s, error: %s.", dir, err)
				return
			}

			outputPath, err = filepath.Abs(outputPath)
			if nil != err {
				logrus.Errorf("get abs output path failed. error: %s.", err)
				return
			}

			title, description, version, contact, apis, err := parser.ParseApisFromComments(dir)
			if nil != err {
				logrus.Errorf("parse failed. error: %s.", err)
				return
			}

			swaggerSpec := parser.NewSwaggerSpec()
			swaggerSpec.Info(
				title,
				description,
				version,
				contact,
			)
			swaggerSpec.Host(host)
			swaggerSpec.Apis(apis)
			swaggerSpec.Schemes([]string{"http", "https"})
			err = swaggerSpec.ParseApis()
			err = swaggerSpec.SaveToFile(outputPath)
			if nil != err {
				logrus.Errorf("save swagger spec to file failed. error: %s.", err)
				return
			}

			return
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&dir, "dir", "d", "./", "search directory")
	flags.StringVarP(&outputPath, "output", "o", "./swagger.json", "swagger json output path")
	flags.StringVarP(&host, "host", "H", "", "host")
	return
}
