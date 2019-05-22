package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"go/importer"

	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	gopath := os.Getenv("GOPATH")
	fmt.Println(gopath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_compile/api/request/api_request.go", nil, parser.AllErrors)
	if nil != err {
		logrus.Errorf("parse file failed. error: %s.", err)
		return
	}

	f2, err := parser.ParseFile(fset, "/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_compile/api/request/model.go", nil, parser.AllErrors)
	if nil != err {
		logrus.Errorf("parse file failed. error: %s.", err)
		return
	}

	conf := types.Config{
		IgnoreFuncBodies: false,
		Importer:         importer.For("source", nil),
		//Importer: importer.For(runtime.Compiler, func(path string) (io.ReadCloser, error) {
		//}),
	}

	info := &types.Info{
		Scopes:     make(map[ast.Node]*types.Scope),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		InitOrder:  make([]*types.Initializer, 0),
	}
	pkg, err := conf.Check("/Users/hao/Documents/Projects/XunLei/video_buddy_service", fset, []*ast.File{f, f2}, info)
	if nil != err {
		logrus.Errorf("check failed. error: %s.", err)
		return
	}

	fmt.Printf("Package  %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())
	fmt.Printf("Imports: %s\n", pkg.Imports())
	fmt.Printf("Scope:   %s\n", pkg.Scope())

	sc := pkg.Scope()
	for _, name := range sc.Names() {
		fmt.Println(sc.Lookup(name))
	}

	//fmt.Printf("Pkg: %#v\n", pkg)
	//fmt.Printf("info.Defs: %#v\n", info.Defs)
	//fmt.Printf("info.Implicits: %#v\n", info.Implicits)
	//fmt.Printf("info.InitOrder: %#v\n", info.InitOrder)
	//fmt.Printf("info.Scopes: %#v\n", info.Scopes)
	//fmt.Printf("info.Selections: %#v\n", info.Selections)
	//fmt.Println("info.Types:")
	//for key, t := range info.Types {
	//	fmt.Printf("-- %s : %#v\n", key, t)
	//}
	//
	//fmt.Printf("info.Uses: %#v\n", info.Uses)
}
