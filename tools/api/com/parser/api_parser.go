package parser

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/sirupsen/logrus"
)

type ApiParser struct {
	Project    *proj.Project
	ProjectDir string
}

func NewApiParser(project *proj.Project) *ApiParser {
	return &ApiParser{
		Project:    project,
		ProjectDir: project.Config.ProjectDir,
	}
}

func (m *ApiParser) ParseRouter() (err error) {
	projectDir := m.ProjectDir

	err = api.NewApiParser(projectDir).MapApi()
	if nil != err {
		logrus.Errorf("mapping api failed. %s.", err)
		return
	}

	return
}
