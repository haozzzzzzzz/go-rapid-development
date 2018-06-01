package init

import (
	"github.com/spf13/cobra"
)

// 初始化API项目
func CommandApiInit() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "api project initialization",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	flags := cmd.Flags()

	return cmd
}
