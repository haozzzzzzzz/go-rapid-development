package service

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/utils/str"
	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateConstant() (err error) {
	serviceDir := m.ServiceDir

	// build constants file
	constantDir := fmt.Sprintf("%s/constant", serviceDir)
	err = os.MkdirAll(constantDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service constant dir %q failed. %s.", constantDir, err)
		return
	}
	constantFilePath := fmt.Sprintf("%s/constant.go", constantDir)
	newConstantFileText := fmt.Sprintf(constantFileText, str.SnakeString(m.Service.Config.Name))
	err = ioutil.WriteFile(constantFilePath, []byte(newConstantFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write file %s failed. error: %s.", constantFilePath, err)
		return
	}
	return
}

// constant.go
var constantFileText = `package constant

const ServiceName = "%s"
`
