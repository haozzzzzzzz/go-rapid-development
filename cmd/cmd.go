package cmd

import (
	"bytes"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 运行命令
func RunCommand(name string, args ...string) (exit int, output *bytes.Buffer, err error) {
	execCommand := exec.Command(name, args...)
	output = bytes.NewBuffer(nil)
	execCommand.Stdout = output
	err = execCommand.Run()
	if nil != err {
		logrus.Errorf("run `%s %s`failed.", name, args)
		exitError, ok := err.(*exec.ExitError)
		if ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			exit = waitStatus.ExitStatus()
		}
		return

	} else {
		waitStatus := execCommand.ProcessState.Sys().(syscall.WaitStatus)
		exit = waitStatus.ExitStatus()
	}

	return
}
