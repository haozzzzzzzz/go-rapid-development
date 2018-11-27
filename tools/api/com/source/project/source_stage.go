package project

import (
	"fmt"
	"io/ioutil"
	"os"

	project2 "github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/utils/str"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func (m *ProjectSource) generateCommonStage() (err error) {
	projectDir := m.ProjectDir
	// stage
	stageDir := fmt.Sprintf("%s/common/stage", projectDir)
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

	envConfig := &struct {
		Debug bool           `json:"debug" yaml:"debug"`
		Stage project2.Stage `json:"stage" yaml:"stage"`
	}{
		Stage: stage,
	}

	logConfig := make(map[string]interface{})
	logConfig["log_level"] = 5
	logConfig["output"] = map[string]interface{}{
		"filedir":     "/Users/hao/Documents/Projects/XunLei/video_buddy_service/bin/log",
		"max_size":    5,
		"max_backups": 3,
		"max_age":     3,
		"compress":    false,
	}

	serviceConfig := &struct {
		MetricsNamespace string `json:"metrics_namespace" yaml:"metrics_namespace"`
	}{
		MetricsNamespace: str.SnakeString(m.Project.Config.Name),
	}

	switch stage {
	case project2.StageDev:
		envConfig.Debug = true

	case project2.StageTest:
		envConfig.Debug = true
		logConfig["output"] = map[string]interface{}{
			"filedir":     "/data/logs",
			"max_size":    500,
			"max_backups": 3,
			"max_age":     3,
			"compress":    false,
		}

	case project2.StagePre, project2.StageProd:
		envConfig.Debug = false
		logConfig["log_level"] = 4
		logConfig["output"] = map[string]interface{}{
			"filedir":     "/data/logs",
			"max_size":    500,
			"max_backups": 3,
			"max_age":     3,
			"compress":    false,
		}

	default:
		err = uerrors.Newf("unknown stage type %s", stage)
		return
	}

	// env.yaml
	envConfigFilePath := fmt.Sprintf("%s/env.yaml", stageConfigDir)
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

	// log.yaml
	logConfigFilePath := fmt.Sprintf("%s/log.yaml", stageConfigDir)
	logConfigFileBytes, err := yaml.Marshal(logConfig)
	if nil != err {
		logrus.Errorf("yaml marshal log config failed. error: %s.", err)
		return
	}

	err = ioutil.WriteFile(logConfigFilePath, logConfigFileBytes, project2.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write log config file failed. path: %s. error: %s.", logConfigFilePath, err)
		return
	}

	// service.yaml
	serviceConfigFilePath := fmt.Sprintf("%s/service.yaml", stageConfigDir)
	serviceConfigFileBytes, err := yaml.Marshal(serviceConfig)
	if nil != err {
		logrus.Errorf("yaml marshal service config failed. error: %s.", err)
		return
	}

	err = ioutil.WriteFile(serviceConfigFilePath, serviceConfigFileBytes, project2.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write service config file failed.path: %s. error: %s.", serviceConfigFilePath, err)
		return
	}

	return
}

var awsConfigFileText = `region: ap-south-1`
var xrayConfigFileText = `daemon_address: "127.0.0.1:3000"
log_level: "warn"
service_version: "1.2.3"`
