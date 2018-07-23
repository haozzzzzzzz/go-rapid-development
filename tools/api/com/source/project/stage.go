package project

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.comozzzzzzzz/go-rapid-development/tools/api/com/project"
	project2 "github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
)

func (m *ProjectSource) generateStage() (err error) {
	projectDir := m.ProjectDir
	// stage
	stageDir := fmt.Sprintf("%s/stage", projectDir)
	err = os.MkdirAll(stageDir, project2.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service stage dir %q failed. %s.", stageDir, err)
		return
	}

	for _, stage := range project2.Stages {
		err = m.generateStageFiles(stageDir, stage)
		if nil != err {
			logrus.Errorf("generate stage %q files failed. %s.", stage, err)
			return
		}
	}

	return
}

func (m *ProjectSource) generateStageFiles(stageDir string, stage project2.Stage) (err error) {
	stageConfigDir := fmt.Sprintf("%s/%s/config", stageDir, stage)
	err = os.MkdirAll(stageConfigDir, project2.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project stage dev dir %q failed. %s.", stageConfigDir, err)
		return
	}

	// env.yaml
	envConfigFilePath := fmt.Sprintf("%s/env.yaml", stageConfigDir)
	envConfig := &struct {
		Debug bool          `json:"debug" yaml:"debug"`
		Stage project.Stage `json:"stage" yaml:"stage"`
	}{
		Stage: stage,
	}

	switch stage {
	case project2.StageDev:
		envConfig.Debug = true
	case project2.StageTest:
		envConfig.Debug = true
	case project2.StagePre:
		envConfig.Debug = false
	case project2.StageProd:
		envConfig.Debug = false
	default:
		err = uerrors.Newf("unknown stage type %s", stage)
		return
	}
	envConfigFileBytes, err := yaml.Marshal(envConfig)
	if nil != err {
		logrus.Errorf("yaml marshal env config failed. %s.", err)
		return
	}

	err = ioutil.WriteFile(envConfigFilePath, envConfigFileBytes, project2.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write env config file %q failed. %s.", envConfigFilePath, err)
		return
	}

	// aws.yaml
	awsConfigFilePath := fmt.Sprintf("%s/aws.yaml", stageConfigDir)
	err = ioutil.WriteFile(awsConfigFilePath, []byte(awsConfigFileText), project2.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write aws config file %q failed. %s.", awsConfigFilePath, err)
		return
	}

	// xray.yaml
	xrayConfigFilePath := fmt.Sprintf("%s/xray.yaml", stageConfigDir)
	err = ioutil.WriteFile(xrayConfigFilePath, []byte(xrayConfigFileText), project2.ProjectFileMode)
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
