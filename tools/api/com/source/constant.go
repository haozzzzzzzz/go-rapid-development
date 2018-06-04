package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateConstant() (err error) {
	projDir := m.ProjectDir

	// build constants file
	constantDir := fmt.Sprintf("%s/constant", projDir)
	err = os.MkdirAll(constantDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project constant dir %q failed. %s.", constantDir, err)
		return
	}
	constantFilePath := fmt.Sprintf("%s/constant.go", constantDir)
	newConstantFileText := fmt.Sprintf(constantFileText, m.Project.Config.Name)
	err = ioutil.WriteFile(constantFilePath, []byte(newConstantFileText), proj.ProjectFileMode)

	return
}

// constant.go
var constantFileText = `package constant

const ServiceName = "%s"
`
