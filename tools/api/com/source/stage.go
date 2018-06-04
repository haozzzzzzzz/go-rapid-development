package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateStage() (err error) {
	projDir := m.ProjectDir
	// stage
	stageDir := fmt.Sprintf("%s/stage", projDir)
	err = os.MkdirAll(stageDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project stage dir %q failed. %s.", stageDir, err)
		return
	}

	for _, stage := range proj.Stages {
		err = m.generateStageFiles(stageDir, stage)
		if nil != err {
			logrus.Errorf("generate stage %q files failed. %s.", stage, err)
			return
		}
	}

	return
}

func (m *ApiProjectSource) generateStageFiles(stageDir string, stage proj.Stage) (err error) {
	stageConfigDir := fmt.Sprintf("%s/%s/config", stageDir, stage)
	err = os.MkdirAll(stageConfigDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project stage dev dir %q failed. %s.", stageConfigDir, err)
		return
	}

	awsConfigFilePath := fmt.Sprintf("%s/aws.yaml", stageConfigDir)
	err = ioutil.WriteFile(awsConfigFilePath, []byte(awsConfigFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write aws config file %q failed. %s.", awsConfigFilePath, err)
		return
	}

	xrayConfigFilePath := fmt.Sprintf("%s/xray.yaml", stageConfigDir)
	err = ioutil.WriteFile(xrayConfigFilePath, []byte(xrayConfigFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write xray config file %q failed. %s.", xrayConfigFilePath, err)
		return
	}

	return
}

var awsConfigFileText = `region: ap-south-1`
var xrayConfigFileText = `daemon_address: "127.0.0.1:3000"
log_level: "warn"
service_version: "1.2.3"`
