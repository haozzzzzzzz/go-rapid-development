package precompile

import (
	"fmt"
	"github.com/spf13/cobra"
)

func CommandPrecompile() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "precompile",
		Short: "precompile go source code",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Println("not implemented")
		},
	}

	return
}
