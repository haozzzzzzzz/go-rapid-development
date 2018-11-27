package project

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func (m *ProjectSource) generateDependent() (err error) {
	serviceDir := m.ProjectDir

	// build dependent file
	dependentDir := fmt.Sprintf("%s/common/dependent", serviceDir)
	err = os.MkdirAll(dependentDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service dependent dir %q failed. error: %s.", dependentDir, err)
		return
	}

	dependentServiceFilePath := fmt.Sprintf("%s/service.go", dependentDir)
	err = ioutil.WriteFile(dependentServiceFilePath, []byte(dependentServiceFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write file %s failed. error: %s.", dependentServiceFilePath, err)
		return
	}

	return
}

// service.go
var dependentServiceFileText = `package dependent

var ServiceName string
`
