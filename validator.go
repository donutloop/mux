package mux

import "net/http"

//Validator validates the incomming value against a valid value/s
type Validator interface {
	Validate(RouteInterface) error
}

//MethodValidator validates the string against a method.
type MethodValidator map[string]struct{}

// newMethodValidator returns default method validator
func newMethodValidator() MethodValidator {
	return MethodValidator(methods)
}

// methods all possible standard methods
var methods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPost:    {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodHead:    {},
	http.MethodPut:     {},
	http.MethodOptions: {},
	http.MethodConnect: {},
}

func (v MethodValidator) Validate(r RouteInterface) error {

	if _, found := v[r.GetMethodName()]; !found {
		return NewBadMethodError(r.GetMethodName())
	}

	return nil
}

//pathValidator check if a path is set and validates the value .
type pathValidator struct{}

func newPathValidator() pathValidator {
	return pathValidator{}
}

func (v pathValidator) Validate(r RouteInterface) error {

	foundPathMatcher := false
	for _, m := range r.GetMatchers() {
		if m.Rank() == rankPath {
			foundPathMatcher = true
		}
	}

	if !foundPathMatcher {
		return NewBadPathError("Patch matcher is missing")
	}

	if len(r.GetPath()) == 0 {
		return NewBadPathError("Path is empty")
	}

	if r.GetPath()[0] != '/' {
		return NewBadPathError("Path starts not with a /")
	}

	return nil
}
