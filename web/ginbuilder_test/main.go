package main

import (
	"log"

	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

func main() {
	parser := ginbuilder.ApiParser{
		ProjectPath: "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/GoLambdaExample/",
		ApiDir:      "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/GoLambdaExample/api",
	}

	err := parser.MapApi()
	if nil != err {
		log.Fatal(err)
		return
	}

}
