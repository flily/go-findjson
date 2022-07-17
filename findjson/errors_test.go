package findjson

import (
	"testing"
)

func TestError(t *testing.T) {
	err := NewJsonError(233, "%s and %s", "foo", "bar")
	if err.Error() != "JSON error at 233: foo and bar" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
