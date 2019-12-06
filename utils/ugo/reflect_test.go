package ugo

import (
	"fmt"
	"testing"
)

func TestStructTags(t *testing.T) {
	o := struct {
		I string `json:" i,omitempty" yaml:"i"`
		J string `json:"j,omitempty"`
		S struct {
		} `json:"s" yaml:"s"`
	}{}
	tags := StructTags(o, []string{"json", "yaml"})

	fmt.Println(tags)

	names := StructJsonTagNames(o)
	fmt.Println(names)
}
