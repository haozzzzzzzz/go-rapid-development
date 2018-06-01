package proj

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ProjectFileMode os.FileMode = os.ModePerm ^ 0111

const ProjectDirMode os.FileMode = os.ModePerm

type ProjectConfigFormat struct {
	Name        string `json:"name" yaml:"name" validate:"required"`
	ProjectDir  string `json:"project_dir" yaml:"project_dir" validate:"required"`
	Description string `json:"description" yaml:"description"`
}

type Project struct {
	Config *ProjectConfigFormat
}

// 加载配置
func (m *Project) Load(projectDir string) (err error) {
	m.Config = &ProjectConfigFormat{}
	confPath := fmt.Sprintf("%s/.proj/proj.yaml", projectDir)
	err = yaml.ReadYamlFromFile(confPath, m.Config)
	if nil != err {
		logrus.Errorf("read project config failed. %s.", err)
		return
	}
	return
}

// 保存配置
func (m *Project) Save() (err error) {
	if m.Config == nil {
		err = errors.New("empty config")
		return
	}

	err = validator.New().Struct(m.Config)
	if nil != err {
		logrus.Errorf("validate project config failed. %s.", err)
		return
	}

	// create project config dir
	projConfDir := fmt.Sprintf("%s/.proj", m.Config.ProjectDir)
	if !file.PathExists(projConfDir) {
		err = os.MkdirAll(projConfDir, ProjectDirMode)
		if nil != err {
			logrus.Errorf("os mkdir %q failed. %s.", projConfDir, err)
			return
		}
	}

	// create project config
	confPath := fmt.Sprintf("%s/.proj/proj.yaml", m.Config.ProjectDir)
	err = yaml.WriteYamlToFile(confPath, m.Config, ProjectFileMode)
	if nil != err {
		logrus.Errorf("save project config to file failed. %s.", err)
		return
	}

	return
}
