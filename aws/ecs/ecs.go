package ecs

import (
	"errors"
	"github.com/haozzzzzzzz/go-rapid-development/utils/ujson"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// https://docs.aws.amazon.com/zh_cn/AmazonECS/latest/developerguide/container-metadata.html
type ContainerMetaData struct {
	MetadataFileStatus string `json:"MetadataFileStatus"`
	PortMappings       []struct {
		HostPort int64 `json:"HostPort"`
	} `json:"PortMappings"`
	HostPrivateIPv4Address string `json:"HostPrivateIPv4Address"`
}

func GetEcsContainerMetaData(retryTimes int, wait time.Duration) (c chan *ContainerMetaData) {
	c = make(chan *ContainerMetaData, 1)
	go func() {
		for i := 0; i < retryTimes; i++ {
			meta, err := getEcsContainerMetaDataFromFile()
			if err != nil {
				logrus.Errorf("get ecs container meta data failed. %#v, error: %s.", meta, err)

			} else {
				c <- meta
				break
			}

			time.Sleep(wait)

		}
	}()

	return
}

func getEcsContainerMetaDataFromFile() (meta *ContainerMetaData, err error) {
	meta = &ContainerMetaData{}

	metaFilePath := os.Getenv("ECS_CONTAINER_METADATA_FILE")
	if metaFilePath == "" {
		err = errors.New("env ECS_CONTAINER_METADATA_FILE is empty")
		return
	}

	err = ujson.ReadJsonFromFile(metaFilePath, meta)
	if nil != err {
		logrus.Errorf("get ecs container meta data from file failed. path: %s, error: %s.", metaFilePath, err)
		return
	}

	if meta.MetadataFileStatus != "READY" {
		err = errors.New("meta data file is not ready")
		return
	}

	return
}
