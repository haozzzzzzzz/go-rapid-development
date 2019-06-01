package book

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

// 更新书本信息
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
			BookName string `json:"book_name" form:"book_name" binding:"required"` //book name
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
