package service

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
)

func (m *ServiceSource) generateBash() (err error) {
	projDir := m.ServiceDir

	shDir := fmt.Sprintf("%s/sh", projDir)
	err = os.MkdirAll(shDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make sh directory %q failed. %s.", shDir, err)
		return
	}

	// params.sh
	err = m.generateBashParams(shDir)
	if nil != err {
		logrus.Errorf("generate params sh failed. %s.", err)
		return
	}

	// build.sh
	err = m.generateBashBuild(shDir)
	if nil != err {
		logrus.Errorf("generate build sh failed. %s.", err)
		return
	}

	// deploy.sh
	err = m.generateBashDeploy(shDir)
	if nil != err {
		logrus.Errorf("generate deploy sh failed. %s.", err)
		return
	}

	// jump_server.sh
	err = m.generateBashJumpServer(shDir)
	if nil != err {
		logrus.Errorf("generate jump server sh failed. %s.", err)
		return
	}

	// target_server.sh
	err = m.generateBashTargetServer(shDir)
	if nil != err {
		logrus.Errorf("generate target server sh failed. %s.", err)
		return
	}

	return
}
