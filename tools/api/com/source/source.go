package source

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

// api 源文件目录
type ApiProjectSource struct {
	Project    *proj.Project
	ProjectDir string
}

func NewApiProjectSource(project *proj.Project) *ApiProjectSource {
	return &ApiProjectSource{
		Project:    project,
		ProjectDir: project.Config.ProjectDir,
	}
}

func (m *ApiProjectSource) Generate() (err error) {

	// generate constant
	err = m.generateConstant()
	if nil != err {
		logrus.Errorf("generate constant failed. %s.", err)
		return
	}

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

	// api
	err = m.generateApi()
	if nil != err {
		logrus.Errorf("generate api failed. %s.", err)
		return
	}

	// main
	err = m.generateMain()
	if nil != err {
		logrus.Errorf("generate main failed. %s.", err)
		return
	}

	// bash
	err = m.generateBash()
	if nil != err {
		logrus.Errorf("generate bash failed. %s.", err)
		return
	}

	return
}
