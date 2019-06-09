package book

import (
	"time"

	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

// 更新书本信息
type ExtraDataItem struct {
	Field1 string `json:"field_1" form:"field_1" binding:"required"`
	Field2 int64  `json:"field_2" form:"field_2" binding:"required"`
}

var BookUpdate ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod: "POST",
	RelativePaths: []string{
		"/api/book/update/:book_id",
	},
	Handle: func(ctx *ginbuilder.Context) (err error) {
		// request path data
		type PathData struct {
			BookId string `json:"book_id" form:"book_id" binding:"required"` // book id
		}
		pathData := PathData{}
		retCode, err := ctx.BindPathData(&pathData)
		if err != nil {
			ctx.Errorf(retCode, "verify update book path data failed. %s.", err)
			return
		}

		// request query data
		type QueryData struct {
			Operator string `json:"operator" form:"operator" binding:"required"` // operator name中午
		}
		queryData := QueryData{}
		retCode, err = ctx.BindQueryData(&queryData)
		if err != nil {
			ctx.Errorf(retCode, "verify update book query data failed. %s.", err)
			return
		}

		// request post data
		type PostData struct {
			BookName       string                   `json:"book_name" form:"book_name" binding:"required"` //book name
			ExtraData      map[string]ExtraDataItem `json:"extra_data" form:"extra_data" binding:"required"`
			Items          []ExtraDataItem          `json:"items" form:"items" binding:"required"`
			InterfaceField interface{}              `json:"interface_field" form:"interface_field" binding:"required"`
		}
		postData := PostData{}
		retCode, err = ctx.BindPostData(&postData)
		if err != nil {
			ctx.Errorf(retCode, "verify update book post data failed. %s.", err)
			return
		}

		// response data
		type ResponseData struct {
			BookId   string `json:"book_id"`   // book id
			BookName string `json:"book_name"` // book name
			Operator string `json:"operator"`  // operator name
		}
		respData := &ResponseData{
			BookId:   pathData.BookId,
			BookName: postData.BookName,
			Operator: queryData.Operator,
		}

		ctx.SuccessReturn(respData)
		return
	},
}

var BookInfo ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod: "GET",
	RelativePaths: []string{
		"/api/book/info/:book_id",
	},
	Handle: func(ctx *ginbuilder.Context) (err error) {
		// request path data
		type PathData struct {
			BookId string `json:"book_id" form:"book_id" binding:"required"` // 书本ID

		}
		pathData := PathData{}
		retCode, err := ctx.BindPathData(&pathData)
		if err != nil {
			ctx.Errorf(retCode, "verify  path data failed. %s.", err)
			return
		}

		// response data
		type ResponseData struct {
			BookId      string `json:"book_id"`      // 书本ID
			BookName    string `json:"book_name"`    // 书本名称
			PublishTime int64  `json:"publish_time"` // 发布时间
		}
		respData := &ResponseData{
			BookId:      pathData.BookId,
			BookName:    "xxx从入门到放弃",
			PublishTime: time.Now().Unix(),
		}

		ctx.SuccessReturn(respData)
		return
	},
}
