package mux

import (
	"fmt"
	"net/http"
)

//Validator validates the incomming value against a valid value/s
type Validator interface {
	Validate(*Route) error
}

//MethodValidator validates the string against a method.
type MethodValidator map[string]struct{}

// newMethodValidator returns default method validator
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

func (v MethodValidator) Validate(r *Route) error {

	if _, found := v[r.methodName]; !found {
		return newBadMethodError(r.methodName)
	}

	return nil
}

//pathMatcherValidator validates the string against a method.
type pathMatcherValidator struct{}

func newPathMatcherValidator() pathMatcherValidator {
	return pathMatcherValidator{}
}

func (v pathMatcherValidator) Validate(r *Route) error {

	for _, m := range r.ms {
		if m.Rank() == rankPath {
			return nil
		}
	}

	return newMissingPathError()
}

//MissingPathError creates error for bad method
type missingPathError struct {
	s string
}

func (bme *missingPathError) Error() string { return fmt.Sprint("Path matcher is missing") }

// newMissingPathErrorreturns an error that formats as the given text.
func newMissingPathError() error {
	return &missingPathError{}
}
