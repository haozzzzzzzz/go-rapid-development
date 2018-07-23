package parser

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/sirupsen/logrus"
)

type ApiParser struct {
	Service    *project.Service
	ServiceDir string
}

func NewApiParser(service *project.Service) *ApiParser {
	return &ApiParser{
		Service:    service,
		ServiceDir: service.Config.ServiceDir,
	}
}

func (m *ApiParser) ParseRouter() (err error) {
	serviceDir := m.ServiceDir

	err = api.NewApiParser(serviceDir).MapApi()
	if nil != err {
		logrus.Errorf("mapping api failed. %s.", err)
		return
	}

	return
}
