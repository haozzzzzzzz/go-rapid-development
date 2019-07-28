package project

import (
	"fmt"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

const ProjectFileMode os.FileMode = os.ModePerm ^ 0111

const ProjectDirMode os.FileMode = os.ModePerm

type ProjectConfigFormat struct {
	Name       string `json:"name" yaml:"name" validate:"required"`
	ProjectDir string `json:"project_dir" yaml:"project_dir" validate:"required"`
}

type Project struct {
	Config *ProjectConfigFormat
}

func (m *Project) Load(projectDir string) (err error) {
	m.Config = &ProjectConfigFormat{}
	confPath := fmt.Sprintf("%s/.project/project.yaml", projectDir)
	err = yaml.ReadYamlFromFile(confPath, m.Config)
	if nil != err {
		logrus.Errorf("read project config %q failed. %s.", confPath, err)
		return
	}
	return
}

func (m *Project) Init() (err error) {
	if m.Config == nil {
		err = errors.New("empty config")
		return
	}

	err = validator.New().Struct(m.Config)
	if nil != err {
		logrus.Errorf("validate service config failed. %s.", err)
		return
	}

	// detect project dir
	projConfDir := fmt.Sprintf("%s/.project/", m.Config.ProjectDir)
	if file.PathExists(projConfDir) {
		err = errors.New("This directory was initialized")
		return
	}

	err = os.MkdirAll(projConfDir, ProjectDirMode)
	if nil != err {
		logrus.Errorf("mkdir %q failed. %s.", projConfDir, err)
		return
	}

	// create project config
	configPath := fmt.Sprintf("%s/project.yaml", projConfDir)
	err = yaml.WriteYamlToFile(configPath, m.Config, ProjectFileMode)
	if nil != err {
		logrus.Errorf("save project config to file failed. %s.", err)
		return
	}

	return
}
