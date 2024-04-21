package skadnetwork

import "errors"

type Validator interface {
	Validate([]byte) (bool, error)
	ValidateString(string) (bool, error)
	Errors() []ValidationError
}

type PostbackValidator struct {
	errors []ValidationError
	Validator
}

// NewPostbackValidator returns a new validator for SKAdNetwork postbacks.
func NewPostbackValidator() Validator {
	return &PostbackValidator{}
}

// Validate validates given postback presented as JSON bytes.
func (v *PostbackValidator) Validate(bytes []byte) (bool, error) {
	v.errors = []ValidationError{}
	p, err := NewPostback(bytes)
	if err != nil {
		return false, err
	}

	ok, err := p.CheckVersion()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("version not supported")
	}

	ok, validationErrors, err := p.ValidateSchema()
	if err != nil {
		return false, err
	}
	if !ok {
		v.errors = append(v.errors, validationErrors...)
		return false, nil
	}

	ok, err = p.VerifySignature()
	if err != nil {
		return false, err
	}
	if !ok {
		v.errors = []ValidationError{NewValidatiorError("invalid signature")}
		return false, nil
	}
	return true, nil
}

// ValidateString validates given postback presented as JSON string.
func (v *PostbackValidator) ValidateString(s string) (bool, error) {
	return v.Validate([]byte(s))
}

// Errors returns all validation errors found by [Validate]. Calling this function does not reset the errors,
// but they will be reset on a subsequent calls to [Validate].
func (v *PostbackValidator) Errors() []ValidationError {
	errors := v.errors
	v.errors = nil

	return errors
}
