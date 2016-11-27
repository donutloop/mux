package mux

import (
	"strings"
	"testing"
)

func TestNewBadRouteError(t *testing.T) {

	r := &Route{
		methodName: "GET",
		path:       "/api/user",
	}

	err := newBadRouteError(r, "Something went wrong")

	if err == nil || !strings.Contains(err.Error(), "Something went wrong") {
		t.Errorf("Bad error message (%v)", err.Error())
	}
}
