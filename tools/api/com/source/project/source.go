package project

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
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
	// stage
	err = m.generateStage()
	if nil != err {
		logrus.Errorf("generate stage failed. %s.", err)
		return
	}

	// com
	err = m.generateCom()
	if nil != err {
		logrus.Errorf("generate com failed. %s.", err)
		return
	}

	return
}