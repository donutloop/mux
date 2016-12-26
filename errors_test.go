package mux

import (
	"strings"
	"testing"
)

func TestBadMethod(t *testing.T) {
	err := NewBadMethodError("GGET")
	if !strings.Contains(err.Error(), "Method not vaild") {
		t.Errorf("Error message is bad (%s)", err.Error())
	}
}

func TestBadPathError(t *testing.T) {
	err := NewBadPathError("/echo")
	if !strings.Contains(err.Error(), "Path is invaild") {
		t.Errorf("Error message is bad (%s)", err.Error())
	}
}

func TestBadRouteError(t *testing.T) {
	err := NewBadRouteError(&Route{}, "Something went wrong")
	if !strings.Contains(err.Error(), "Something went wrong") {
		t.Errorf("Error message is bad (%s)", err.Error())
	}
}
