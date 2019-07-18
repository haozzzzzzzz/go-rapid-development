package parser

import (
	"os"

	"go/parser"
	"go/token"

	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
)

func ParseApisFromComments(dir string) (
	title string,
	description string,
	version string,
	contact string,
	apis []*ApiItem,
	err error,
) {
	title = "api service name"
	description = "api service description"
	version = "1.0"
	contact = "contact"
	apis = make([]*ApiItem, 0)

	pkgDirs := make([]string, 0)
	pkgDirs, err = file.SearchFileNames(dir, func(fileInfo os.FileInfo) bool {
		if fileInfo.IsDir() {
			return true
		} else {
			return false
		}
	}, true)
	pkgDirs = append(pkgDirs, dir)

	for _, pkgDir := range pkgDirs {
		subApis, errParse := ParseApisFromPkgComment(pkgDir)
		err = errParse
		if nil != err {
			logrus.Errorf("parse apis from pkg comment failed. pkgDir: %s, error: %s.", pkgDir, err)
			return
		}

		apis = append(apis, subApis...)
	}

	return
}

func ParseApisFromPkgComment(pkgDir string) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, pkgDir, nil, parser.ParseComments)
	if nil != err {
		logrus.Errorf("parse pkg dir failed. pkgDir: %s, error: %s.", pkgDir, err)
		return
	}

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Files {
			for _, commentGroup := range astFile.Comments {
				for _, comment := range commentGroup.List {
					tempApis, errParse := ParseApisFromPkgCommentText(comment.Text)
					err = errParse
					if nil != err {
						logrus.Errorf("parse apis from pkg comment text failed. text: %s, error: %s.", comment.Text, err)
						return
					}

					apis = append(apis, tempApis...)
				}
			}
		}
	}
	return
}

func ParseApisFromPkgCommentText(text string) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)

	// TODO
	return
}
