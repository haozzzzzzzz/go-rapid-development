package cmd

import (
	"context"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 运行命令
func RunCommand(
	workingDir string,
	name string,
	args ...string,
) (exit int, err error) {
	execCommand := exec.Command(name, args...)
	execCommand.Dir = workingDir
	output, err := execCommand.Output()
	strOutput := string(output)
	if strOutput != "" {
		logrus.Info(strOutput)
	}

	if nil != err {
		logrus.Errorf("run `%s %s`failed.", name, args)
		exitError, ok := err.(*exec.ExitError)
		if ok {
			err = exitError
			logrus.Error(string(exitError.Stderr))
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

func RunCommandCtx(
	ctx context.Context,
	workingDir string,
	name string,
	args ...string,
) (exit int, err error) {
	execCommand := exec.CommandContext(ctx, name, args...)
	execCommand.Dir = workingDir
	output, err := execCommand.Output()
	strOutput := string(output)
	if strOutput != "" {
		logrus.Info(strOutput)
	}

	if nil != err {
		logrus.Errorf("run `%s %s`failed.", name, args)
		exitError, ok := err.(*exec.ExitError)
		if ok {
			err = exitError
			logrus.Error(string(exitError.Stderr))
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
