package skadnetwork

import (
	"github.com/xeipuuv/gojsonschema"
)

type SchemaHelper struct {
	*gojsonschema.Schema
}

func NewSchemaHelper(version string) (SchemaHelper, error) {
	var helper SchemaHelper
	schemaLoader := gojsonschema.NewReferenceLoader("https://raw.githubusercontent.com/whisk/skadnetwork/main/schema/v" + version)
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
