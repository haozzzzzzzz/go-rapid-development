package source

import (
	"fmt"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateCom() (err error) {
	projDir := m.ProjectDir

	comDir := fmt.Sprintf("%s/com", projDir)
	err = os.MkdirAll(comDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project dir %s failed. %s.", comDir, err)
		return
	}

	// config
	err = m.generateComConfig(comDir)
	if nil != err {
		logrus.Errorf("generate com config failed. %s.", err)
		return
	}

	// metrics
	err = m.generateComMetrics(comDir)
	if nil != err {
		logrus.Errorf("generate com metrics failed. %s.", err)
		return
	}

	return
}
