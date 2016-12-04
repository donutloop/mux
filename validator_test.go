package mux

import (
	"fmt"
	"strings"
	"testing"
)

func TestMethodValidator(t *testing.T) {
	validator := newMethodValidator()

	for k := range methods {
		t.Run(fmt.Sprintf("Method: %s", k), func(t *testing.T) {
			r := &Route{
				methodName: k,
			}

			if err := validator.Validate(r); err != nil {
				t.Errorf("Unexpected invalid method (%s)", k)
			}
		})
	}
}

func TestMethodValidatorFail(t *testing.T) {
	validator := newMethodValidator()
	r := &Route{
		methodName: "GGET",
	}
	if err := validator.Validate(r); err == nil || !strings.Contains(err.Error(), "Method not vaild") {
		t.Errorf("Unexpected valid method (%s)", r.methodName)
	}
}
