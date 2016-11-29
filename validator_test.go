package mux

import (
	"fmt"
	"testing"
)

func TestMethodValidator(t *testing.T) {
	validator := newMethodValidator()
	for k := range methods {
		t.Run(fmt.Sprintf("Method: %s", k), func(t *testing.T) {
			if err := validator.Validate(k); err != nil {
				t.Errorf("Unexpected method (%s)", k)
			}
		})
	}
}

func TestMethodValidatorFail(t *testing.T) {
	validator := newMethodValidator()
	method := "GGET"
	if err := validator.Validate(method); err == nil {
		t.Errorf("Unexpected valid method (%s)", method)
	}
}
