#!/usr/bin/env bash
export GOROOT=/usr/local/go
export GOOS=darwin

# api project
go build -o api ../../tools/api/main.go

## log formatter
#go build -o logfmt ${GOPATH}/src/github.com/haozzzzzzzz/go-rapid-development/tools/logfmt/main.go
#
## lambda build
#go build -o lamb ${GOPATH}/src/github.com/haozzzzzzzz/go-lambda/tools/lamb/main.go
#
## lambda deploy
#go build -o lamd ${GOPATH}/src/github.com/haozzzzzzzz/go-lambda/tools/lamd/main.go

echo finish