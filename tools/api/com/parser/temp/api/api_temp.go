package api

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

type Request struct {
	ReqKey1 string `json:"req_key_1" form:"req_key_1"`
}

var TestApi ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod:   "GET",
	RelativePath: "/test/:name",
	Handle: func(ctx *ginbuilder.Context) (err error) {
		// request path data
		type PathData struct {
			Request
			PathKey1 string `json:"path_key_1" form:"path_key_1"`
			PathKey2 int    `json:"path_key_2" form:"path_key_2"`
		}
		pathData := PathData{}
		retCode, err := ctx.BindPathData(&pathData)
		if err != nil {
			ctx.Errorf(retCode, "verify  path data failed. %s.", err)
			return
		}

		// request query data
		type QueryData struct {
			QueryKey1 string `json:"query_key_1" form:"query_key_1"`
			QueryKey2 string `json:"query_key_2" form:"query_key_2"`
		}
		queryData := &QueryData{}
		retCode, err = ctx.BindQueryData(queryData)
		if err != nil {
			ctx.Errorf(retCode, "verify  query data failed. %s.", err)
			return
		}

		// request post data
		type PostData struct {
			Request
			PostKey1 string `json:"post_key_1" form:"post_key_1"`
			PostKey2 int64  `json:"post_key_2" form:"post_key_2"`
		}
		postData := PostData{}
		retCode, err = ctx.BindPostData(&postData)
		if err != nil {
			ctx.Errorf(retCode, "verify  post data failed. %s.", err)
			return
		}

		// response data
		type ResponseData struct {
			RsKey string `json:"rs_key" form:"rs_key"`
		}
		respData := &ResponseData{}

		ctx.SuccessReturn(respData)
		return
	},
}
