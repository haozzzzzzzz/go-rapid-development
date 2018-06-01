package proj

import (
	"fmt"
	"testing"
)

func TestProject_Save(t *testing.T) {
	config := &ProjectConfigFormat{
		Name:       "test",
		ProjectDir: "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/common/proj",
	}

	project := &Project{
		Config: config,
	}

	err := project.Save()
	if nil != err {
		t.Error(err)
		return
	}
}

func TestProject_Load(t *testing.T) {
	project := &Project{}
	err := project.Load("/Users/hao/Documents/Projects/Github/go_lambda_learning/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/common/proj")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(project.Config)
}
