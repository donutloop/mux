package mux

import (
	"net/http"
	"testing"
)

func TestGetQueriesFail(t *testing.T) {
	r := &http.Request{}

	if value := GetQueries(r); value != nil {
		t.Errorf("Unexpected value (%v)", value)
	}
}

func TestCurrentRouteFail(t *testing.T) {
	r := &http.Request{}

	if value := CurrentRoute(r); value != nil {
		t.Errorf("Unexpected value (%v)", value)
	}
}

func TestGetVarsFail(t *testing.T) {
	r := &http.Request{}

	if value := GetVars(r); value != nil {
		t.Errorf("Unexpected value (%v)", value)
	}
}
