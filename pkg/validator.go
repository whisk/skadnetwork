package skadnetwork // import "github.com/whisk/skadnetwork/pkg"

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
	p, err := NewPostback(bytes)
	if err != nil {
		return false, err
	}
	ok, validatationErrors, err := p.ValidateSchema()
	if err != nil {
		return false, err
	}
	if !ok {
		v.errors = validatationErrors
		return false, nil
	}

	ok, err = p.VerifySignature()
	return ok, err
}

func (v *PostbackValidator) ValidateString(s string) (bool, error) {
	return v.Validate([]byte(s))
}

func (v *PostbackValidator) Errors() []ValidationError {
	errors := v.errors
	v.errors = nil

	return errors
}
