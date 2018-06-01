package cmd

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/compile"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/doc"
	_init "github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/init"
)

func init() {
	rootCmd.AddCommand(_init.CommandApiInit())
	rootCmd.AddCommand(compile.CommandApiCompile())
	rootCmd.AddCommand(doc.CommandApiDoc())
}
