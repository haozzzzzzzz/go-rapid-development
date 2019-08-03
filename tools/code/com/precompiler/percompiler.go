/**
go source precompiler package
please work with go +build constraint
*/
package precompiler

import (
	"fmt"
	"github.com/gosexy/to"
	"github.com/haozzzzzzzz/go-rapid-development/tools/lib/gofmt"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"strings"
)

func PrecompileText(
	filename string,
	src interface{},
	params map[string]interface{},
) (newFileText string, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.ParseComments)
	if nil != err {
		logrus.Errorf("parser parse file failed. error: %s.", err)
		return
	}

	// pos是从1开始的
	type Block struct {
		Text       string
		StartIndex int
		EndIndex   int // 不包含EndIndex所在的Index
	}

	preStatements := make([]*Block, 0)

	for _, commentGroup := range f.Comments {
		for _, comment := range commentGroup.List {
			strStatement := strings.Replace(comment.Text, "//", "", 1)
			strStatement = strings.TrimSpace(strStatement)
			if !strings.HasPrefix(strStatement, "+pre") {
				continue
			}
			strStatement = strings.Replace(strStatement, "+pre", "", 1)
			strStatement = strings.TrimSpace(strStatement)

			statement := &Block{
				Text:       strStatement,
				StartIndex: int(comment.Pos()) - 1,
				EndIndex:   int(comment.End()),
			}

			preStatements = append(preStatements, statement)
		}
	}

	bFileText, err := ioutil.ReadFile(filename)
	if nil != err {
		logrus.Errorf("read file failed. error: %s.", err)
		return
	}
	lenOfFile := len(bFileText)
	_ = lenOfFile

	goBlocks := make([]*Block, 0)
	luaBlocks := make([]*Block, 0)

	var lastStatement *Block
	var blocksStart, blocksEnd int
	_ = blocksEnd
	for _, statement := range preStatements {
		if lastStatement == nil {
			blocksStart = statement.StartIndex
			lastStatement = statement

		} else {
			blockText := string(bFileText[lastStatement.EndIndex:statement.StartIndex])
			block := &Block{
				Text:       blockText,
				StartIndex: lastStatement.EndIndex,
				EndIndex:   statement.StartIndex,
			}

			idxGoBlocks := len(goBlocks)
			goBlocks = append(goBlocks, block)

			luaBlock := &Block{
				StartIndex: block.StartIndex,
				EndIndex:   block.EndIndex,
			}
			luaBlocks = append(luaBlocks, luaBlock)

			luaBlock.Text = fmt.Sprintf("use_go_block(%d)\n", idxGoBlocks)

			lastStatement = statement
		}

		blocksEnd = lastStatement.EndIndex
		luaBlocks = append(luaBlocks, statement)
	}

	var luaText string
	for _, block := range luaBlocks {
		luaText += block.Text + "\n"
	}

	Lua := lua.NewState()
	defer Lua.Close()

	useIndexes := make([]int, 0)
	Lua.SetGlobal("use_go_block", Lua.NewFunction(func(state *lua.LState) int {
		lv := state.ToInt(1) // get first argument
		useIndexes = append(useIndexes, lv)
		return 0
	}))

	for key, param := range params {
		var val lua.LValue
		kind := reflect.TypeOf(param).Kind()
		switch {
		case kind == reflect.Bool:
			val = lua.LBool(param.(bool))
		case (kind >= reflect.Int && kind <= reflect.Uint64) || (kind >= reflect.Float32 && kind <= reflect.Float64):
			val = lua.LNumber(to.Int64(param))
		case kind == reflect.String:
			val = lua.LString(to.String(param))
		default:
			err = uerrors.Newf("unsupported go type to lua type. %s", kind)
		}

		Lua.SetGlobal(key, val)
	}

	err = Lua.DoString(luaText)
	if nil != err {
		logrus.Errorf("lua do string failed. error: %s.", err)
		return
	}

	useGoMap := make(map[int]bool)
	for _, useIdx := range useIndexes {
		useGoMap[useIdx] = true
	}

	bNewFileText := make([]byte, 0)

	if blocksStart > 0 {
		bNewFileText = bFileText[:blocksStart]
	}

	for idx, goBlock := range goBlocks {
		if useGoMap[idx] { // 包含
			bNewFileText = append(bNewFileText, bFileText[goBlock.StartIndex:goBlock.EndIndex]...)
		}
	}

	if len(bNewFileText) < blocksEnd {
		bNewFileText = append(bNewFileText, bFileText[blocksEnd:]...)
	}

	newFileText = string(bNewFileText)
	newFileText, err = gofmt.StrGoFmt(newFileText)
	if nil != err {
		logrus.Errorf("fmt go src failed. error: %s.", err)
		return
	}

	return
}
