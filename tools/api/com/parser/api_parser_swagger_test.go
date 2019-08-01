package parser

import (
	"fmt"
	"testing"
)

func TestSaveApisSwaggerSpec(t *testing.T) {
	swgSpc := NewSwaggerSpec()
	swgSpc.Apis([]*ApiItem{
		{
			PackageName:    "pack",
			ApiHandlerFunc: "func",
			HttpMethod:     "GET",
			RelativePaths: []string{
				"/api/book/:book_id",
				"/api/book",
			},
			PathData: &StructType{
				Fields: []*Field{
					{
						TypeName: "int64",
						Tags: map[string]string{
							"json":    "book_id",
							"binding": "required",
						},
					},
				},
			},
		},
	})
	swgSpc.Info(
		"Book shop",
		"book shop api for testing tools",
		"1",
		"haozzzzzzzz",
	)
	err := swgSpc.ParseApis()
	if nil != err {
		t.Error(err)
		return
	}

	out, err := swgSpc.Output()
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(string(out))
}
