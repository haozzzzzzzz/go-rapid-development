package cmd

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/code/cmd/ertn"
	"github.com/haozzzzzzzz/go-rapid-development/tools/code/cmd/precompile"
)

var RootCmd = rootCmd

func init() {
	RootCmd.AddCommand(precompile.CommandPrecompile())
	RootCmd.AddCommand(ertn.CommandErtn())
}
