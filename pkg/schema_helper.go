package skadnetwork

import (
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

type SchemaHelper struct {
	*gojsonschema.Schema
}

func NewSchemaHelper(version string) (SchemaHelper, error) {
	var helper SchemaHelper
	schemaPath, err := filepath.Abs("schema/v" + version)
	if err != nil {
		return helper, err
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return helper, err
	}

	helper.Schema = schema
	return helper, nil
}

func (s SchemaHelper) Validate(p Postback) (bool, []ValidationError, error) {
	res, err := s.Schema.Validate(gojsonschema.NewBytesLoader(p.bytes))
	if err != nil {
		return false, nil, err
	}
	if res.Valid() {
		return true, nil, nil
	}

	errors := []ValidationError{}
	for _, e := range res.Errors() {
		errors = append(errors, NewValidatiorError(e.String()))
	}
	return false, errors, nil
}
