package request

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

type PostDataField struct {
	Field1 string `json:"filed_1" form:"filed_1"`
	Field2 struct {
		SubField1 string `json:"sub_field_1" form:"sub_field_1"`
	} `json:"field_2" form:"field_2"`
}

var TestRequest ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod:   "POST",
	RelativePath: "/test_request",
	Handle: func(ctx *ginbuilder.Context) (err error) {
		// request post data
		type PostData struct {
			//ArrayType []struct {
			//	SubArrayType []string `json:"sub_array_type" form:"sub_array_type"`
			//} `json:"array_type" form:"array_type"`
			//Inter          interface{}         `json:"inter" form:"inter"`
			//PostDataField1 *PostDataField      `json:"post_data_field_1" form:"post_data_field_1"`
			//PostDataField2 []*PostDataField    `json:"post_data_field_2" form:"post_data_field_2"`
			//PostDataField3 interface{}         `json:"post_data_field_3" form:"post_data_field_3"`
			//PostDataField4 []interface{}       `json:"post_data_field_4" form:"post_data_field_4"`
			PostDataField5 *OtherFilePostData1 `json:"post_data_field_5" form:"post_data_field_5"`
			PostDataField6 string              `json:"post_data_field_6" form:"post_data_field_6"`
		}
		postData := PostData{}
		retCode, err := ctx.BindPostData(&postData)
		if err != nil {
			ctx.Errorf(retCode, "verify  post data failed. %s.", err)
			return
		}

		ctx.Success()
		return
	},
}
