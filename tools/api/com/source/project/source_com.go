package project

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
)

func (m *ProjectSource) generateCommonCom() (err error) {
	projectDir := m.ProjectDir

	comDir := fmt.Sprintf("%s/common/com", projectDir)
	err = os.MkdirAll(comDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service component dir %s failed. %s.", comDir, err)
		return
	}

	// config
	err = m.generateCommonComConfig(comDir)
	if nil != err {
		logrus.Errorf("generate com config failed. %s.", err)
		return
	}

	// metrics
	err = m.generateCommonComMetrics(comDir)
	if nil != err {
		logrus.Errorf("generate com metrics failed. %s.", err)
		return
	}

	return
}
