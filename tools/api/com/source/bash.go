package source

import (
	"fmt"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateBash() (err error) {
	projDir := m.ProjectDir

	shDir := fmt.Sprintf("%s/sh", projDir)
	err = os.MkdirAll(shDir, proj.ProjectDirMode)
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
