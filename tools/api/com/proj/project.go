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

func LoadProject(projectDir string) (project *Project, err error) {
	project = &Project{}
	err = project.Load(projectDir)
	if nil != err {
		logrus.Errorf("project load config failed. %s.", err)
		return
	}

	return
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

// 初始化
func (m *Project) Init() (err error) {
	if m.Config == nil {
		err = errors.New("empty config")
		return
	}

	err = validator.New().Struct(m.Config)
	if nil != err {
		logrus.Errorf("validate project config failed. %s.", err)
		return
	}

	// detect project dir
	projectDir := m.Config.ProjectDir
	if !file.PathExists(projectDir) {
		err = os.MkdirAll(projectDir, ProjectDirMode)
		if nil != err {
			logrus.Errorf("make project directory %q failed. %s.", projectDir, err)
			return
		}
	}

	// create project config dir
	projConfDir := fmt.Sprintf("%s/.proj", projectDir)
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