package imports

import (
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
)

func CommandApiImports() (cmd *cobra.Command) {
	var path string
	cmd = &cobra.Command{
		Use: "imports",
		Short: "api dependency imports",
		Run: func(cmd *cobra.Command, args []string) {
			if path == "" {
				logrus.Errorf("empty imports path")
				return
			}

			// goimports
			goimports.DoGoImports([]string{path}, true)

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&path, "path", "p", "./", "imports path")

	return
}
