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
