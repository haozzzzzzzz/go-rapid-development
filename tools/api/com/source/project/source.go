package project

import (
	"fmt"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

type ProjectSource struct {
	Project    *project.Project
	ProjectDir string
}

func NewProjectSource(project *project.Project) *ProjectSource {
	return &ProjectSource{
		Project:    project,
		ProjectDir: project.Config.ProjectDir,
	}
}

func (m *ProjectSource) Generate() (err error) {
	projectDir := m.ProjectDir

	// config
	err = m.generateCommonConfig()
	if nil != err {
		logrus.Errorf("generate config failed. error: %s.", err)
		return
	}

	// stage
	err = m.generateCommonStage()
	if nil != err {
		logrus.Errorf("generate stage failed. %s.", err)
		return
	}

	// metrics
	err = m.generateCommonMetrics()
	if nil != err {
		logrus.Errorf("generate metrics failed. %s.", err)
		return
	}

	// generate app service root dir
	appDir := fmt.Sprintf("%s/app", projectDir)
	err = os.MkdirAll(appDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project service dir %s failed. %s.", appDir, err)
		return
	}

	// generate manage service root dir
	manageDir := fmt.Sprintf("%s/manage", projectDir)
	err = os.MkdirAll(manageDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project service dir %s failed. %s.", manageDir, err)
		return
	}

	// generate rpc service root dir
	rpcDir := fmt.Sprintf("%s/rpc", projectDir)
	err = os.MkdirAll(rpcDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project rpc service dir failed. dir: %s.error: %s.", projectDir, err)
		return
	}

	return
}
