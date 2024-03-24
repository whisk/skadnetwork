package skadnetwork // import "github.com/whisk/skadnetwork/pkg"

import (
	"encoding/json"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

type Validator interface {
	Validate([]byte) (bool, error)
	ValidateString(string) (bool, error)
	Errors() []ValidationError
}

type PostbackValidator struct {
	errors []ValidationError
	Validator
}

func NewPostbackValidator() Validator {
	return &PostbackValidator{}
}

func (v *PostbackValidator) Validate(bytes []byte) (bool, error) {
	ok, err := v.ValidateSchema(bytes)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	ok, err = v.Verify(bytes)
	return ok, err
}

func (v *PostbackValidator) ValidateString(s string) (bool, error) {
	return v.Validate([]byte(s))
}

func (v *PostbackValidator) ValidateSchema(bytes []byte) (bool, error) {
	schemaPath, err := filepath.Abs("schema/v4.0")
	if err != nil {
		return false, err
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return false, err
	}
	res, err := schema.Validate(gojsonschema.NewBytesLoader(bytes))
	if err != nil {
		return false, err
	}
	if res.Valid() {
		return true, nil
	}
	v.errors = []ValidationError{}
	for _, e := range res.Errors() {
		v.errors = append(v.errors, NewValidatiorError(e.String()))
	}
	return false, nil
}

func (v *PostbackValidator) Verify(bytes []byte) (bool, error) {
	var p Postback
	err := json.Unmarshal(bytes, &p)
	if err != nil {
		return false, err
	}
	return p.Verify()
}

func (v *PostbackValidator) VerifyString(s string) (bool, error) {
	return v.Verify([]byte(s))
}

func (v *PostbackValidator) Errors() []ValidationError {
	errors := v.errors
	v.errors = nil

	return errors
}
