package cmd

import (
	"testing"
)

func TestRunCommand(t *testing.T) {
	_, err := RunCommand("echo", "hello, world")
	if nil != err {
		t.Error(err)
		return
	}
}
