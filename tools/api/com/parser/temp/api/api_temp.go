package api

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

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
