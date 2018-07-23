package service

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
)

// api 源文件目录
type ServiceSource struct {
	Service    *project.Service
	ServiceDir string
}

func NewServiceSource(service *project.Service) *ServiceSource {
	return &ServiceSource{
		Service:    service,
		ServiceDir: service.Config.ServiceDir,
	}
}

type GenerateParams struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (m *ServiceSource) Generate(params *GenerateParams) (err error) {

	// generate constant
	err = m.generateConstant()
	if nil != err {
		logrus.Errorf("generate constant failed. %s.", err)
		return
	}

	// api
	err = m.generateApi()
	if nil != err {
		logrus.Errorf("generate api failed. %s.", err)
		return
	}

	// main
	err = m.generateMain(params)
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