package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
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
		Importer: importer.Default(),
	}

	pkg, err := conf.Check("/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_compile/api/request", fset, []*ast.File{f, f2}, nil)
	if nil != err {
		logrus.Errorf("check failed. error: %s.", err)
		return
	}

	fmt.Printf("Package  %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())
	fmt.Printf("Imports: %s\n", pkg.Imports())
	fmt.Printf("Scope:   %s\n", pkg.Scope())
	fmt.Printf("Pkg: %#v\n", pkg)
}
