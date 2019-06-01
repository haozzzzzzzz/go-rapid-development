package lswagger

import (
	"github.com/go-openapi/spec"
)

type Swagger struct {
	spec.Swagger
}

func NewSwagger() (swg *Swagger) {
	swg = &Swagger{}
	swg.Paths = NewPaths()
	swg.Swagger.Swagger = "2.0"
	return
}

func (m *Swagger) PathsAdd(key string, pathItem *spec.PathItem) {
	m.Paths.Paths[key] = *pathItem
}

func NewPaths() *spec.Paths {
	return &spec.Paths{
		Paths: make(map[string]spec.PathItem),
	}
}
