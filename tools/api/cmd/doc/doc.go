package doc

import (
	"github.com/spf13/cobra"
)

func CommandApiDoc() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "doc",
		Short: "api doc",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}
