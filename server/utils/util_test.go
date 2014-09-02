package utils

import (
	"errors"
	"testing"
)

func TestRandomString(t *testing.T) {
	n := 7
	s := RandomString(n)
	if len(s) != n {
		t.Error(errors.New("Invalid random string size"))
	}
}
