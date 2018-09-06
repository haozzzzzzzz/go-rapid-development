package project

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func (m *ProjectSource) generateCommonConfig() (err error) {
	configDir := fmt.Sprintf("%s/common/config", m.ProjectDir)
	err = os.MkdirAll(configDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make com config dir %q failed. %s.", configDir, err)
		return
	}

	// init
	initFilePath := fmt.Sprintf("%s/init.go", configDir)
	err = ioutil.WriteFile(initFilePath, []byte(initFileText), project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("write init file failed. file: %s. error: %s.", initFilePath, err)
		return
	}

	// env
	envFilePath := fmt.Sprintf("%s/env.go", configDir)
	err = ioutil.WriteFile(envFilePath, []byte(envFileText), project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("write env file %q failed. %s.", envFilePath, err)
		return
	}

	// aws
	awsFilePath := fmt.Sprintf("%s/aws.go", configDir)
	err = ioutil.WriteFile(awsFilePath, []byte(awsFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write aws file %q failed. %s.", awsFilePath, err)
		return
	}

	// xray
	xrayFilePath := fmt.Sprintf("%s/xray.go", configDir)
	err = ioutil.WriteFile(xrayFilePath, []byte(xrayFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write xray file %q  failed. %s.", xrayFilePath, err)
		return
	}

	// log
	logFilePath := fmt.Sprintf("%s/log.go", configDir)
	err = ioutil.WriteFile(logFilePath, []byte(logFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write log file failed. path: %s. error: %s.", logFilePath, err)
		return
	}

	// service
	serviceFilePath := fmt.Sprintf("%s/service.go", configDir)
	err = ioutil.WriteFile(logFilePath, []byte(serviceFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write service file failed. path: %s. error: %s.", serviceFilePath, err)
		return
	}

	return
}

var initFileText = `package config

import (
	"log"

	"github.com/sirupsen/logrus"
)

func init() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	newLogger := logrus.New()
	newLogger.Formatter = &logrus.JSONFormatter{}
	log.SetOutput(newLogger.Writer())

	CheckAWSConfig()
	CheckEnvConfig()
	CheckServiceConfig()
	CheckXRayConfig()
}
`

var envFileText = `package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

type EnvConfigFormat struct {
	Debug bool   ` + "`json:\"debug\" yaml:\"debug\"`" + `
	Stage string ` + "`json:\"stage\" yaml:\"stage\" validate:\"required\"`" + `
}

func (m *EnvConfigFormat) WithStagePrefix(strVal string) string {
	return fmt.Sprintf("%s_%s", m.Stage, strVal)
}

var EnvConfig *EnvConfigFormat

func CheckEnvConfig() {
	if EnvConfig != nil {
		return
	}

	EnvConfig = &EnvConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/env.yaml", EnvConfig)
	if nil != err {
		logrus.Errorf("read env config file failed. %s.", err)
		return
	}

	err = validator.New().Struct(EnvConfig)
	if nil != err {
		logrus.Fatalf("validate env config failed. %s", err)
		return
	}

}
`

var awsFileText = `package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/aws/ec2"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

type AWSConfigFormat struct {
	Region string ` + "`yaml:\"region\" validate:\"required\"`" + `
}

var AWSConfig *AWSConfigFormat
var AWSSession *session.Session
var AWSEc2InstanceIdentifyDocument ec2metadata.EC2InstanceIdentityDocument

func CheckAWSConfig() {
	if AWSConfig != nil {
		return
	}

	AWSConfig = &AWSConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/aws.yaml", AWSConfig)
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

	AWSEc2InstanceIdentifyDocument, err = ec2.GetEc2InstanceIdentityDocument(AWSSession)
	if nil != err {
		logrus.Errorf("get ec2 instance identify document failed. %s.", err)
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

var XRayConfig *XRayConfigFormat

func CheckXRayConfig() {
	if XRayConfig != nil {
		return
	}

	XRayConfig = &XRayConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/xray.yaml", XRayConfig)
	if nil != err {
		logrus.Fatalf("read xray config failed. %s", err)
		return
	}

	err = validator.New().Struct(XRayConfig)
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

var logFileText = `package config

import (
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

type LogConfigFormat struct {
	LogLevel logrus.Level ` + "`json:\"log_level\" yaml:\"log_level\" validate:\"required\"`" + `
}

var LogConfig *LogConfigFormat

func CheckLogConfig() {
	if LogConfig != nil {
		return
	}

	LogConfig = &LogConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/log.yaml", LogConfig)
	if nil != err {
		logrus.Fatalf("read log config file failed. %s", err)
		return
	}

	err = validator.New().Struct(LogConfig)
	if nil != err {
		logrus.Fatalf("validate log config failed. %s", err)
		return
	}

	logrus.SetLevel(LogConfig.LogLevel)
}
`

var serviceFileText = `package config

import (
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

type ServiceConfigFormat struct {
	MetricsNamespace string ` + "`json:\"metrics_namespace\" yaml:\"metrics_namespace\" validate:\"required\"`" + `
}

var ServiceConfig *ServiceConfigFormat

func CheckServiceConfig() {
	if ServiceConfig != nil {
		return
	}

	ServiceConfig = &ServiceConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/service.yaml", ServiceConfig)
	if nil != err {
		logrus.Fatalf("read service config file failed. %s", err)
		return
	}

	err = validator.New().Struct(ServiceConfig)
	if nil != err {
		logrus.Errorf("validate service config failed. %s.", err)
		return
	}

}
`
