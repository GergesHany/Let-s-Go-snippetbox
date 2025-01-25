package assert

import (
	"strings"
	"testing"
)

// Equal asserts that two values are equal.
func Equal[T comparable](t *testing.T, actual, expected T) {	
	t.Helper() // report the filename and line number of the code which called our Equal() function in the output.
    if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
	   t.Errorf("The error from StringContains")
	//    t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
	   t.Errorf("got: %v; expected: nil", actual)
	}
}