package request

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_compile/api/request/model"
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
		// request path data
		type PathData struct {
			Day string `json:"day" form:"day"`
		}
		pathData := PathData{}
		retCode, err := ctx.BindPathData(&pathData)
		if err != nil {
			ctx.Errorf(retCode, "verify  path data failed. %s.", err)
			return
		}

		// request query data
		type QueryData struct {
			Key string `json:"key" form:"key"`
		}
		queryData := QueryData{}
		retCode, err = ctx.BindQueryData(&queryData)
		if err != nil {
			ctx.Errorf(retCode, "verify  query data failed. %s.", err)
			return
		}

		// request post data
		type PostData struct {
			ArrayType []struct {
				SubArrayType []string `json:"sub_array_type" form:"sub_array_type"`
			} `json:"arrays_type" form:"array_type"`
			//Inter          interface{}         `json:"inter" form:"inter"`
			PostDataField1 *PostDataField      `json:"post_data_field_1" form:"post_data_field_1"`
			PostDataField2 []*PostDataField    `json:"post_data_field_2" form:"post_data_field_2"`
			PostDataField3 interface{}         `json:"post_data_field_3" form:"post_data_field_3"`
			PostDataField4 []interface{}       `json:"post_data_field_4" form:"post_data_field_4"`
			PostDataField5 *OtherFilePostData1 `json:"post_data_field_5" form:"post_data_field_5"`
			PostDataField6 string              `json:"post_data_field_6" form:"post_data_field_6"`
			PostDataField7 model.Model2        `json:"post_data_field_7" form:"post_data_field_7"`
			PostDataField8 struct {
				SubField1 string `json:"sub_field_1" form:"sub_field_1"`
			} `json:"post_data_field_8" form:"post_data_field_8"`
			PostDataField9  map[string]string        `json:"post_data_field_9" form:"post_data_field_9"`
			PostDataField10 map[string]*model.Model2 `json:"post_data_field_10" form:"post_data_field_10"`
		}
		postData := PostData{}
		retCode, err = ctx.BindPostData(&postData)
		if err != nil {
			ctx.Errorf(retCode, "verify  post data failed. %s.", err)
			return
		}

		// response data
		type ResponseData struct {
			Hello string `json:"hello"`
		}
		respData := &ResponseData{}

		ctx.SuccessReturn(respData)
		return
	},
}

var TestPaths ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod: "GET",
	RelativePaths: []string{
		"/api/test_paths_1",
		"/api/test_paths_2",
	},
	Handle: func(ctx *ginbuilder.Context) (err error) {
		ctx.Success()
		return
	},
}
