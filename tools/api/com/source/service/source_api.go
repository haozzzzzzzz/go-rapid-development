package service

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateApi() (err error) {
	serviceDir := m.ServiceDir

	// generate api folder
	apiDir := fmt.Sprintf("%s/api", serviceDir)
	err = os.MkdirAll(apiDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make service api dir %q failed. %s.", apiDir, err)
		return
	}

	// routers
	routersFilePath := fmt.Sprintf("%s/routers.go", apiDir)
	err = ioutil.WriteFile(routersFilePath, []byte(routersFileText), project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("write api/routers.go failed. \n%s.", err)
		return
	}

	return
}

// routers.go
var routersFileText = `package api

import (
	"github.com/gin-gonic/gin"
)

// 注意：BindRouters函数体内不能自定义添加任何声明，由api compile命令生成api绑定声明
func BindRouters(engine *gin.Engine) (err error) {
	return
}
`
