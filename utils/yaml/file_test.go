package yaml

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestReadYamlFromFile(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	_ = ctx

	type TestStruct struct {
		TestUint32 logrus.Level `yaml:"log_level"`
	}

	config := &TestStruct{}
	err = ReadYamlFromFile("/Users/hao/Documents/Projects/XunLei/video_buddy_user_device/common/stage/dev/config/consul/log.yaml", config)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println("val: ",config.TestUint32)
}
