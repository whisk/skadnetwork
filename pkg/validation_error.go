package skadnetwork

type ValidationError struct {
	message string
}

func NewValidatiorError(msg string) ValidationError {
	return ValidationError{message: msg}
}

func (e *ValidationError) Error() string {
	return e.message
}
