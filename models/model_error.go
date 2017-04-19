package models

// ModelError : Specific model error
type ModelError struct {
	s string
	c string
}

var (
	// InvalidInputCode : provided input is not valid
	InvalidInputCode = "M0001" // => ErrBadReqBody
	// InternalCode : internal error
	InternalCode = "M0002" // => ErrInternal
	// TimeoutCode : A timeout on microservice communication happened
	TimeoutCode = "M0003" // => ErrGatewayTimeout
)

// NewError : Returns new model error
func NewError(c, s string) *ModelError {
	return &ModelError{c: c, s: s}
}

// Error : Returns the error string, and implements go error interface
func (e *ModelError) Error() string {
	return e.s
}

// Code : Returns the error code
func (e *ModelError) Code() string {
	return e.c
}
