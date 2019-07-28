package project

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type ServiceType string

const (
	ServiceTypeApp    ServiceType = "app"
	ServiceTypeManage ServiceType = "manage"
	ServiceTypeRPC    ServiceType = "rpc"
)

type ServiceConfigFormat struct {
	Name        string `json:"name" yaml:"name" validate:"required"`
	ServiceDir  string `json:"service_dir" yaml:"service_dir" validate:"required"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"type" yaml:"type" validate:"required"`
}

type Service struct {
	Config *ServiceConfigFormat
}

func LoadService(serviceDir string) (service *Service, err error) {
	service = &Service{}
	err = service.Load(serviceDir)
	if nil != err {
		logrus.Errorf("service load config failed. %s.", err)
		return
	}

	return
}

// 加载配置
func (m *Service) Load(serviceDir string) (err error) {
	m.Config = &ServiceConfigFormat{}
	confPath := fmt.Sprintf("%s/.service/service.yaml", serviceDir)
	err = yaml.ReadYamlFromFile(confPath, m.Config)
	if nil != err {
		logrus.Errorf("read service config failed. %s.", err)
		return
	}
	return
}

// 初始化
func (m *Service) Init() (err error) {
	if m.Config == nil {
		err = errors.New("empty config")
		return
	}

	m.Config.ServiceDir, err = filepath.Abs(m.Config.ServiceDir)
	if nil != err {
		logrus.Errorf("get service dir failed. error: %s.", err)
		return
	}

	err = validator.New().Struct(m.Config)
	if nil != err {
		logrus.Errorf("validate service config failed. %s.", err)
		return
	}

	// detect service dir
	serviceDir := m.Config.ServiceDir

	// create service config dir
	srvConfDir := fmt.Sprintf("%s/.service", serviceDir)
	if file.PathExists(srvConfDir) {
		err = errors.New("This directory was initialized")
		return
	}

	err = os.MkdirAll(srvConfDir, ProjectDirMode)
	if nil != err {
		logrus.Errorf("os mkdir %q failed. %s.", srvConfDir, err)
		return
	}

	// create service config
	confPath := fmt.Sprintf("%s/.service/service.yaml", m.Config.ServiceDir)
	err = yaml.WriteYamlToFile(confPath, m.Config, ProjectFileMode)
	if nil != err {
		logrus.Errorf("save service config to file failed. %s.", err)
		return
	}

	return
}
