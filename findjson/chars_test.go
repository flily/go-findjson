package findjson

import (
	"bytes"
	"testing"
)

func testIsXXXInCharSet(t *testing.T, name string, set []byte, f func(byte) bool) {
	for _, c := range set {
		if !f(c) {
			t.Errorf("'%02x' is not %s", c, name)
		}
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		if bytes.IndexByte(set, c) < 0 && f(c) {
			t.Errorf("'%02x' is %s", c, name)
		}
	}
}

func TestIsWhiteSpace(t *testing.T) {
	set := []byte{
		' ', '\t', '\n', '\r',
	}

	testIsXXXInCharSet(t, "whitespace", set, isWhiteSpace)
}

func TestIsEscapeChar(t *testing.T) {
	set := []byte{
		'"', '\\', '/', 'b', 'f', 'n', 'r', 't',
	}

	testIsXXXInCharSet(t, "escape char", set, isEscapeChar)
}

func TestIsDigit(t *testing.T) {
	set := []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	testIsXXXInCharSet(t, "digit", set, isDigit)
}

func TestIsNonZeroDigit(t *testing.T) {
	set := []byte{
		'1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	testIsXXXInCharSet(t, "non-zero digit", set, isNonZeroDigit)
}

func TestIsHexDigits(t *testing.T) {
	set := []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	}

	testIsXXXInCharSet(t, "hex digit", set, isHexDigit)
}
