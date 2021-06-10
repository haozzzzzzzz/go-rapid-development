package umapstructure

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func Decode(in interface{}, out interface{}, tag string) (err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   out,
		TagName:  tag,
	})
	if err != nil {
		logrus.Errorf("new mapstructure decoder failed. error: %s", err)
		return
	}

	err = decoder.Decode(in)
	if err != nil {
		logrus.Errorf("mapstructure decode failed. error: %s", err)
		return
	}
	return
}
