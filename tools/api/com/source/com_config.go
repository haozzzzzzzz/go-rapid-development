package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateComConfig(comDir string) (err error) {
	configDir := fmt.Sprintf("%s/config", comDir)
	err = os.MkdirAll(configDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make com config dir %q failed. %s.", configDir, err)
		return
	}

	// aws
	awsFilePath := fmt.Sprintf("%s/aws.go", configDir)
	err = ioutil.WriteFile(awsFilePath, []byte(awsFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write aws file %q failed. %s.", awsFilePath, err)
		return
	}

	// xray
	xrayFilePath := fmt.Sprintf("%s/xray.go", configDir)
	err = ioutil.WriteFile(xrayFilePath, []byte(xrayFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write xray file %q  failed. %s.", xrayFilePath, err)
		return
	}

	return
}

var awsFileText = `package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

type AWSConfigFormat struct {
	Region string ` + "`yaml:\"region\" validate:\"required\"`" + `
}

var AWSConfig AWSConfigFormat
var AWSSession *session.Session

func init() {
	var err error
	err = yaml.ReadYamlFromFile("./config/aws.yaml", &AWSConfig)
	if nil != err {
		logrus.Fatalf("read aws config file failed. %s", err)
		return
	}

	err = validator.New().Struct(AWSConfig)
	if nil != err {
		logrus.Fatalf("validate aws config file failed. %s", err)
		return
	}

	awsConfig := &aws.Config{}
	awsConfig.Region = aws.String(AWSConfig.Region)
	AWSSession, err = session.NewSession(awsConfig)
	if nil != err {
		logrus.Fatalf("new aws session failed. %s", err)
		return
	}

}
`

var xrayFileText = `package config

import (
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-xray-sdk-go/xray"
	// Importing the plugins enables collection of AWS resource information at runtime.
	// Every plugin should be imported after "github.com/aws/aws-xray-sdk-go/xray" library.
	_ "github.com/aws/aws-xray-sdk-go/plugins/ec2"
)

type XRayConfigFormat struct {
	DaemonAddress  string ` + "`json:\"daemon_address\" yaml:\"daemon_address\" validate:\"required\"`" + `
	LogLevel       string ` + "`json:\"log_level\" yaml:\"log_level\"`" + `
	ServiceVersion string ` + "`json:\"service_version\" yaml:\"service_version\"`" + `
}

var XRayConfig XRayConfigFormat

func init() {
	var err error
	err = yaml.ReadYamlFromFile("./config/xray.yaml", &XRayConfig)
	if nil != err {
		logrus.Fatalf("read xray config failed. %s", err)
		return
	}

	err = validator.New().Struct(&XRayConfig)
	if nil != err {
		logrus.Fatalf("validate xray config failed. %s", err)
		return
	}

	err = xray.Configure(xray.Config{
		DaemonAddr:     XRayConfig.DaemonAddress, 
		LogLevel:       XRayConfig.LogLevel,
		ServiceVersion: XRayConfig.ServiceVersion,
	})
	if nil != err {
		logrus.Fatalf("configure xray config failed. %s", err)
		return
	}
}
`
