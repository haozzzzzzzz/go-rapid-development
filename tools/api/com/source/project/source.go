package project

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
)

type ProjectSource struct {
	Project *project.Project
	ProjectDir string
}

func NewProjectSource(project *project.Project) *ProjectSource {
	return &ProjectSource{
		Project: project,
		ProjectDir: project.Config.ProjectDir,
	}
}

func (m *ProjectSource) Generate() (err error) {
	projectDir := m.ProjectDir

	// stage
	err = m.generateCommonStage()
	if nil != err {
		logrus.Errorf("generate stage failed. %s.", err)
		return
	}

	// com
	err = m.generateCommonCom()
	if nil != err {
		logrus.Errorf("generate com failed. %s.", err)
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
	err = os.MkdirAll(manageDir, project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("make project service dir %s failed. %s.", manageDir, err)
		return
	}

	return
}