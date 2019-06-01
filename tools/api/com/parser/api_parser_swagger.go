package parser

import (
	"strings"

	"fmt"

	"io/ioutil"

	"github.com/go-openapi/spec"
	"github.com/haozzzzzzzz/go-rapid-development/api/request"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/tools/lib/lswagger"
	"github.com/sirupsen/logrus"
)

type SwaggerSpec struct {
	apis    []*ApiItem
	Swagger *lswagger.Swagger
}

func NewSwaggerSpec() (swgSpec *SwaggerSpec) {
	swgSpec = &SwaggerSpec{
		apis:    make([]*ApiItem, 0),
		Swagger: lswagger.NewSwagger(),
	}
	return
}

func (m *SwaggerSpec) ParseApis() (
	err error,
) {
	logrus.Info("Save swagger spec ...")
	defer func() {
		if nil != err {
			logrus.Errorf("Save swagger spec failed. %s", err)
			return
		} else {
			logrus.Info("Save swagger spec finish")
		}
	}()

	swag := lswagger.NewSwagger()
	for _, api := range m.apis {
		paths := api.RelativePaths
		for _, path := range paths {
			// 将gin url path上的:变量转换成swagger的{变量}
			subPaths := strings.Split(path, "/")
			for i, subPath := range subPaths {
				if strings.HasPrefix(subPath, ":") {
					subPath = strings.Replace(subPath, ":", "", 1)
					subPath = fmt.Sprintf("{%s}", subPath)
					subPaths[i] = subPath
				}
			}

			path = strings.Join(subPaths, "/")
			err = m.parseApi(path, api)
			if nil != err {
				logrus.Errorf("swagger spec parse api failed. error: %s.", err)
				return
			}

		}
	}

	_ = swag
	return
}

func (m *SwaggerSpec) parseApi(path string, api *ApiItem) (err error) {
	pathItem := &spec.PathItem{}
	operation := &spec.Operation{}
	operation.Consumes = []string{request.MIME_JSON}
	operation.Produces = []string{request.MIME_JSON}
	operation.Summary = api.Summary
	operation.Description = api.Description
	operation.ID = fmt.Sprintf("%s-%s", api.PackageFuncName(), path)
	operation.Parameters = make([]spec.Parameter, 0)

	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#pathsObject
	switch api.HttpMethod {
	case request.METHOD_GET:
		pathItem.Get = operation
	case request.METHOD_POST:
		pathItem.Post = operation
	case request.METHOD_PUT:
		pathItem.Put = operation
	case request.METHOD_DELETE:
		pathItem.Delete = operation
	case request.METHOD_OPTIONS:
		pathItem.Options = operation
	case request.METHOD_HEAD:
		pathItem.Head = operation
	case request.METHOD_PATCH:
		pathItem.Patch = operation
	default:
		logrus.Warnf("not supported method for swagger spec. method: %s", api.HttpMethod)
	}

	// path data
	if api.PathData != nil {
		for _, pathField := range api.PathData.Fields {
			operation.Parameters = append(operation.Parameters, *FieldBasicParameter("path", pathField))
		}
	}

	// query data
	if api.QueryData != nil {
		for _, queryField := range api.QueryData.Fields {
			operation.Parameters = append(operation.Parameters, *FieldBasicParameter("query", queryField))
		}
	}

	// post data
	if api.PostData != nil {
		body := &spec.Parameter{}
		body.In = "body"
		body.Name = api.PostData.Name
		body.Description = api.PostData.Description
		body.Required = true
		body.Schema = ITypeToSwaggerSchema(api.PostData)
		operation.Parameters = append(operation.Parameters, *body)
	}

	operation.Responses = &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			StatusCodeResponses: map[int]spec.Response{
				200: {
					ResponseProps: spec.ResponseProps{
						Description: "success",
					},
				},
			},
		},
	}

	m.Swagger.PathsAdd(path, pathItem)

	return
}

func (m *SwaggerSpec) Apis(apis []*ApiItem) {
	m.apis = apis
}

func (m *SwaggerSpec) Host(host string) {
	m.Swagger.Host = host
}

func (m *SwaggerSpec) Schemes(schemes []string) {
	m.Swagger.Schemes = schemes
}

func (m *SwaggerSpec) Info(
	title string,
	description string,
	version string,
	contactName string,
) {
	m.Swagger.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       title,
			Description: description,
			Version:     version,
			Contact: &spec.ContactInfo{
				Name: contactName,
			},
		},
	}

	return
}

func (m *SwaggerSpec) SaveToFile(fileName string) (err error) {
	out, err := m.Output()
	if nil != err {
		logrus.Errorf("get spec output failed. error: %s.", err)
		return
	}

	err = ioutil.WriteFile(fileName, out, project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("save spec to file failed. error: %s.", err)
		return
	}

	return
}

func (m *SwaggerSpec) Output() (output []byte, err error) {
	return m.Swagger.MarshalJSON()
}

// query、path基础类型参数
func FieldBasicParameter(in string, field *Field) (parameter *spec.Parameter) {
	parameter = &spec.Parameter{}
	parameter.Name = field.TagJson()
	parameter.In = in
	parameter.Description = field.Description
	parameter.Required = field.Required()
	switch field.TypeSpec.(type) {
	case *BasicType:
		parameter.Type = BasicTypeTransformSchemaType(field.TypeName)

	default:
		parameter.Type = BasicTypeTransformSchemaType(field.TypeName)

	}

	return
}

func BasicTypeTransformSchemaType(fieldType string) (swagType string) {
	switch fieldType {
	case "string":
		swagType = "string"

	case "bool":
		swagType = "boolean"

	default:
		if strings.Contains(fieldType, "float") {
			swagType = "number"
		} else {
			swagType = "integer"
		}
	}
	return
}

func ITypeToSwaggerSchema(iType IType) (schema *spec.Schema) {
	schema = &spec.Schema{}
	switch iType.(type) {
	case *StructType:
		structType := iType.(*StructType)
		schema.Type = []string{"object"}
		schema.Required = make([]string, 0)
		schema.Properties = make(map[string]spec.Schema)

		for _, field := range structType.Fields {
			jsonName := field.TagJson()
			fieldSchema := ITypeToSwaggerSchema(field.TypeSpec)
			fieldSchema.Description = field.Description
			if field.Required() {
				fieldSchema.Required = []string{jsonName}
			}

			schema.Properties[jsonName] = *fieldSchema
		}

	case *BasicType:
		basicType := iType.(*BasicType)
		schemaType := BasicTypeTransformSchemaType(basicType.Name)
		schema.Type = []string{schemaType}
	}

	return
}
