package cmd

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/compile"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/doc"
	_init "github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/init"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/cmd/imports"
)

func init() {
	rootCmd.AddCommand(_init.CommandApiInit())
	rootCmd.AddCommand(_init.CommandApiAddService())
	rootCmd.AddCommand(compile.CommandApiCompile())
	rootCmd.AddCommand(imports.CommandApiImports())
	rootCmd.AddCommand(doc.CommandApiDoc())
}
