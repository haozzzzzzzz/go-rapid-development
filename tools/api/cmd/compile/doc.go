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
	var importSource bool
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
			apis, err := apiParser.ScanApis(true, importSource)
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
	flags.BoolVarP(&importSource, "import_source", "s", false, "import source")
	return cmd
}

func GenerateCommentDoc() (cmd *cobra.Command) {
	var dir string
	cmd = &cobra.Command{
		Use:   "com_doc",
		Short: "generate open api doc from comment",
		Run: func(cmd *cobra.Command, args []string) {
			return
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&dir, "dir", "d", "./", "search directory")
	return
}
