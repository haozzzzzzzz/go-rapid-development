package service

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateFirstInit() (err error) {
	serviceDir := m.ServiceDir

	// first init folder
	firstInitFolder := fmt.Sprintf("%s/first_init", serviceDir)
	err = os.MkdirAll(firstInitFolder, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service first init dir %s failed. error: %s.", firstInitFolder, err)
		return
	}

	firstInitFilePath := fmt.Sprintf("%s/init.go", firstInitFolder)
	err = ioutil.WriteFile(firstInitFilePath, []byte(firstInitFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write file %s failed. error: %s.", firstInitFilePath, err)
		return
	}

	return
}

var firstInitFileText = `package first_init

func init() {
	dependent.ServiceName = constant.ServiceName
}
`
