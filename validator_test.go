package mux

import (
	"fmt"
	"strings"
	"testing"
)

func TestBadMethod(t *testing.T) {
	err := newBadMethodError("GGET")
	if err.Error() == "Method not vaild" {
		t.Errorf("Error message is bad (%s)", err.Error())
	}
}

func TestMethodValidator(t *testing.T) {
	validator := newMethodValidator()
	for k := range methods {
		t.Run(fmt.Sprintf("Method: %s", k), func(t *testing.T) {
			if err := validator.Validate(k); err != nil {
				t.Errorf("Unexpected invalid method (%s)", k)
			}
		})
	}
}

func TestMethodValidatorFail(t *testing.T) {
	validator := newMethodValidator()
	method := "GGET"
	if err := validator.Validate(method); err == nil || !strings.Contains(err.Error(), "Method not vaild") {
		t.Errorf("Unexpected valid method (%s)", method)
	}
}
