package api

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

type Request struct {
	ReqKey1 string `json:"req_key_1" form:"req_key_1"`
}

type Field struct {
	//FieldKey1 string      `json:"field_key_1" form:"field_key_1"`
	FieldKey2 FieldChild1 `json:"field_key_2" form:"field_key_2"`
}

type FieldChild1 struct {
	//FieldChild1Key1 string      `json:"field_child_1_key_1" form:"field_child_1_key_1"`
	FiledChild1Key2 FiledChild2 `json:"filed_child_1_key_2" form:"filed_child_1_key_2"`
}

type FiledChild2 struct {
	FiledChild2Key string `json:"filed_child_2_key" form:"filed_child_2_key"`
}

var TestApi ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod:   "GET",
	RelativePath: "/test/:name",
	Handle: func(ctx *ginbuilder.Context) (err error) {
		// request post data
		type PostData struct {
			//Request
			//PostKey1      string         `json:"post_key_1" form:"post_key_1"`
			//PostKey2      int64          `json:"post_key_2" form:"post_key_2"`
			PostFieldKey3 *Field `json:"post_field_key_3" form:"post_field_key_3"`
			//PostMapKey4   map[int]*Field `json:"post_map_key_4" form:"post_map_key_4"`
		}
		postData := PostData{}
		retCode, err := ctx.BindPostData(&postData)
		if err != nil {
			ctx.Errorf(retCode, "verify  post data failed. %s.", err)
			return
		}

		return
	},
}
