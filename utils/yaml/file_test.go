package yaml

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/haozzzzzzzz/go-rapid-development/consul"
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

	type ConsulConfigFormat struct {
		ClientConfig *consul.ClientConfigFormat `yaml:"client_config" validate:"required,required"`
		KeyPrefix    string                     `yaml:"key_prefix" validate:"required"`
	}

	config := &ConsulConfigFormat{}
	err = ReadYamlFromFile("/Users/hao/Documents/Projects/XunLei/video_buddy_user_device/common/stage/dev/config/consul.yaml", config)
	if nil != err {
		logrus.Errorf("read yaml from file failed. error: %s.", err)
		return
	}

	fmt.Println(config)
}
