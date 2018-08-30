package cmd

import (
	"testing"
)

func TestRunCommand(t *testing.T) {
	_, err := RunCommand("/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/service/VideoBuddyHomeApi/lambda/LambdaPlayfailSnsHomeSubscriber/stage/dev", "/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/service/VideoBuddyHomeApi/lambda/LambdaPlayfailSnsHomeSubscriber/detector", "hello, world")
	if nil != err {
		t.Error(err)
		return
	}
}
