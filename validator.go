package mux

import (
	"fmt"
	"net/http"
)

//Validator validates the incomming value against a valid value/s
type Validator interface {
	Validate(string) error
}

//MethodValidator matches the string against a method.
type MethodValidator map[string]struct{}

// newMatcher
func newMethodValidator() MethodValidator {
	return MethodValidator(methods)
}

//BadMethodError creates error for bad method
type badMethodError struct {
	s string
}

func (bme *badMethodError) Error() string { return fmt.Sprintf("Method not vaild (%s)", bme.s) }

// newBadMethodError returns an error that formats as the given text.
func newBadMethodError(text string) error {
	return &badMethodError{text}
}

// methods all possible standard methods
var methods = map[string]struct{}{
	http.MethodGet:     struct{}{},
	http.MethodPost:    struct{}{},
	http.MethodPatch:   struct{}{},
	http.MethodDelete:  struct{}{},
	http.MethodHead:    struct{}{},
	http.MethodPut:     struct{}{},
	http.MethodOptions: struct{}{},
	http.MethodConnect: struct{}{},
}

func (m MethodValidator) Validate(method string) error {

	if _, found := m[method]; !found {
		return newBadMethodError(method)
	}

	return nil
}
