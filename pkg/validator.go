package skadnetwork

import (
	"errors"
)

type PostbackValidator struct {
	errors []ValidationError
}

// NewPostbackValidator returns a new validator for SKAdNetwork postbacks.
func NewPostbackValidator() PostbackValidator {
	return PostbackValidator{}
}

// Validate performs all validations for a given postback presented as JSON bytes:
//  - JSON schema validation according to the postback version
//  - signature verification
// Validate returns the validation result and an error. Non-nil error indicates that the validation itself
// has failed and we are not sure if the postback is valid or not.
func (v *PostbackValidator) Validate(bytes []byte) (bool, error) {
	v.errors = []ValidationError{}
	p, err := NewPostback(bytes)
	if err != nil {
		// error when initializing the postback means it is invalid
		v.errors = append(v.errors, NewValidationError(err.Error()))
		return false, nil
	}

	ok, err := p.VersionSupported()
	if err != nil {
		// error when checking for the version means the postback is invalid
		v.errors = append(v.errors, NewValidationError(err.Error()))
		return false, nil
	}
	if !ok {
		// unsupported version means we can't say if it is valid or not
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
		v.errors = []ValidationError{NewValidationError("invalid signature")}
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
