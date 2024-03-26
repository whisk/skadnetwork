package skadnetwork

import "fmt"

type ValidationError struct {
	message string
	fmt.Stringer
	error
}

func NewValidatiorError(msg string) ValidationError {
	return ValidationError{message: msg}
}

func (e *ValidationError) String() string {
	return e.message
}
