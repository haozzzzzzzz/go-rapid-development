/**
check error return info
*/
package ertn

import (
	"fmt"
	"github.com/spf13/cobra"
)

func CommandErrorReturn() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "ertn",
		Short: "check error return",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Println("not implemented")
		},
	}

	return
}
