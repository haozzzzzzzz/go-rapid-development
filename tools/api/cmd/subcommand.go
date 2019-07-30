package cmd

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/compile"
	_init "github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/init"
)

func init() {
	rootCmd.AddCommand(_init.CommandApiInit())
	rootCmd.AddCommand(_init.CommandApiAddService())
	rootCmd.AddCommand(compile.CommandApiCompile())
	rootCmd.AddCommand(compile.GenerateApiDoc())
	rootCmd.AddCommand(compile.GenerateCommentDoc())
}
