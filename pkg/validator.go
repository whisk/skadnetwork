package skadnetwork

import (
	"errors"
)

var supportedVersions = map[string]bool{"2.1": true, "2.2": true, "3.0": true, "4.0": true}

type PostbackValidator struct {
	p *Postback
	errors []ValidationError
}

// NewPostbackValidator returns a new validator for SKAdNetwork postbacks.
func NewPostbackValidator() PostbackValidator {
	return PostbackValidator{}
}

// Validate performs all validations for a given postback presented as JSON bytes:
//  - checks if the postback version is supported by the validator
//  - runs all checks
// Validate returns the validation result and an error. Non-nil error indicates that the validation itself
// has failed and we are not sure if the postback is valid or not.
func (v *PostbackValidator) Validate(bytes []byte) (bool, error) {
	err := v.Init(bytes)
	if err != nil {
		// error when initializing the postback means it is invalid
		v.errors = append(v.errors, NewValidationError(err.Error()))
	}
	return v.Check()
}

// Check performs all actual checks: schema validation and signature verification
func (v *PostbackValidator) Check() (bool, error) {
	if v.p == nil {
		return false, errors.New("validator not initialized")
	}
	ok, err := v.VersionSupported()
	if err != nil {
		// error when checking for the version means the postback is invalid
		v.errors = append(v.errors, NewValidationError(err.Error()))
		return false, nil
	}
	if !ok {
		// unsupported version means we can't say if it is valid or not
		return false, errors.New("version not supported")
	}

	ok, validationErrors, err := v.p.ValidateSchema()
	if err != nil {
		return false, err
	}
	if !ok {
		v.errors = append(v.errors, validationErrors...)
		return false, nil
	}

	ok, err = v.p.VerifySignature()
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

func (v *PostbackValidator) Init(bytes []byte) error{
	v.errors = []ValidationError{}
	v.p = nil
	p, err := NewPostback(bytes)
	if err != nil {
		return err
	}
	v.p = &p
	return nil
}

func (v *PostbackValidator) VersionSupported() (bool, error) {
	if v.p == nil {
		return false, errors.New("validator not initialized")
	}
	version, ok := v.p.Version()
	if !ok {
		return false, errors.New("empty or invalid version value")
	}
	_, ok = supportedVersions[version]
	return ok, nil

}

// Errors returns all validation errors found by [Validate]. Calling this function does not reset the errors,
// but they will be reset on a subsequent calls to [Validate].
func (v *PostbackValidator) Errors() []ValidationError {
	errors := v.errors
	v.errors = nil

	return errors
}
