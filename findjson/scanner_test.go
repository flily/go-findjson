package findjson

import (
	"fmt"
	"strings"
	"testing"
)

func TestBufferStartsWith(t *testing.T) {
	//                0        10        20        30        40
	buffer := []byte("the quick brown fox jumps over the lazy dog")

	{
		s := []byte("the")
		r, l := bufferStartsWith(buffer, 0, s)
		if !r || l != 3 {
			t.Errorf("bufferStartsWith(buffer, 0, '%s') returns %v, %d", string(s), r, l)
		}
	}

	{
		s := []byte("thank")
		r, l := bufferStartsWith(buffer, 0, s)
		if r || l != 2 {
			t.Errorf("bufferStartsWith(buffer, 0, '%s') returns %v, %d", string(s), r, l)
		}
	}

	{
		s := []byte("THE")
		r, l := bufferStartsWith(buffer, 0, s)
		if r || l != 0 {
			t.Errorf("bufferStartsWith(buffer, 0, '%s') returns %v, %d", string(s), r, l)
		}
	}

	{
		s := []byte("quick")
		r, l := bufferStartsWith(buffer, 4, s)
		if !r || l != 5 {
			t.Errorf("bufferStartsWith(buffer, 4, '%s') returns %v, %d", string(s), r, l)
		}
	}

	{
		s := []byte("QUICK")
		r, l := bufferStartsWith(buffer, 4, s)
		if r || l != 0 {
			t.Errorf("bufferStartsWith(buffer, 4, '%s') returns %v, %d", string(s), r, l)
		}
	}

	{
		s := []byte("dog")
		r, l := bufferStartsWith(buffer, 40, s)
		if !r || l != 3 {
			t.Errorf("bufferStartsWith(buffer, 40, '%s') returns %v, %d", string(s), r, l)
		}
	}

	{
		s := []byte("doggie")
		r, l := bufferStartsWith(buffer, 40, s)
		if r || l != 3 {
			t.Errorf("bufferStartsWith(buffer, 40, '%s') returns %v, %d", string(s), r, l)
		}
	}
}

func TestBufferFindSample(t *testing.T) {
	//                0        10        20        30        40
	buffer := []byte("the quick brown fox jumps over the lazy dog")

	if s := bufferFindSample(buffer, 0, 10); s != "the quick " {
		t.Errorf("bufferFindSample(buffer, 0, 10) returns '%s'", s)
	}

	if s := bufferFindSample(buffer, 40, 10); s != "dog" {
		t.Errorf("bufferFindSample(buffer, 40, 10) returns '%s'", s)
	}

	if s := bufferFindSample(buffer, 50, 0); s != "EOF" {
		t.Errorf("bufferFindSample(buffer, 50, 0) returns '%s'", s)
	}
}

func TestJumpNextNonWhiteSpace(t *testing.T) {
	//                0        10        20        30        40
	buffer := []byte("the quick             brown     fox        ")

	if i := jumpNextNonWhiteSpace(buffer, 0); i != 0 {
		t.Errorf("jumpNextNonWhiteSpace(buffer, 0) returns %d", i)
	}

	if i := jumpNextNonWhiteSpace(buffer, 10); i != 22 {
		t.Errorf("jumpNextNonWhiteSpace(buffer, 10) returns %d", i)
	}

	if i := jumpNextNonWhiteSpace(buffer, 40); i != len(buffer) {
		t.Errorf("jumpNextNonWhiteSpace(buffer, 40) returns %d", i)
	}

	if i := jumpNextNonWhiteSpace(buffer, 50); i != len(buffer) {
		t.Errorf("jumpNextNonWhiteSpace(buffer, 50) returns %d", i)
	}
}

type scannerCorrectCases []string

func (c scannerCorrectCases) On(t *testing.T, s JsonTokenScanner) {
	for _, item := range c {
		info := fmt.Sprintf("run scanner '%+v'", s)
		t.Run(info, func(tt *testing.T) {
			b := []byte(item)
			start, end, err := s(b, 0)
			if err != nil {
				tt.Errorf("%v returns %d, %d, %s", s, start, end, err)
				sample := bufferFindSample(b, err.(*JsonError).Offset, 20)
				tt.Errorf("  sample: %s", sample)
				tt.Errorf("  case: %s", item)
				tt.Errorf("        %s^", strings.Repeat(" ", end))
			}

			if start != 0 || end != len(b) {
				tt.Errorf("%v returns %d, %d, nil", s, start, end)
			}
		})
	}
}

func TestScanJsonLiteralSuccess(t *testing.T) {
	caseList := scannerCorrectCases{
		"true",
		"false",
		"null",
	}

	caseList.On(t, scanJsonLiteral)
}

func TestScanJsonLiteralFailure(t *testing.T) {
	{
		//           0         1         2
		//           012345678901234567890
		s := []byte(` true false null tr`)
		//                           |^
		start, end, err := scanJsonLiteral(s, 17)
		if err == nil {
			t.Fatalf("scanJsonLiteral(s, 17) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 17: expect null, true or false, got 'tr'" {
			t.Errorf("scanJsonLiteral(s, 17) returns %d, %d, %s", start, end, err)
		}

		if start != 17 || end != 19 {
			t.Errorf("scanJsonLiteral(s, 17) returns %d, %d, nil", start, end)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte(` true false null tr`)
		//                            |^
		start, end, err := scanJsonLiteral(s, 18)
		if err == nil {
			t.Fatalf("scanJsonLiteral(s, 18) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 18: expect null, true or false, got 'r'" {
			t.Errorf("scanJsonLiteral(s, 18) returns %d, %d, %s", start, end, err)
		}

		if start != 18 || end != 18 {
			t.Errorf("scanJsonLiteral(s, 18) returns %d, %d, nil", start, end)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte(` true false null tr`)
		//                             |^
		start, end, err := scanJsonLiteral(s, 19)
		if err == nil {
			t.Fatalf("scanJsonLiteral(s, 19) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 19: expect null, true or false, got 'EOF'" {
			t.Errorf("scanJsonLiteral(s, 19) returns %d, %d, %s", start, end, err)
		}

		if start != 19 || end != 19 {
			t.Errorf("scanJsonLiteral(s, 19) returns %d, %d, nil", start, end)
		}
	}
}

func TestScanDigits(t *testing.T) {
	{
		//           0         1         2
		//           012345678901234567890
		s := []byte("123456abcdefghijkl")
		//           |<->|
		start, end, err := scanDigits(s, 0)
		if err != nil {
			t.Errorf("scanDigits(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 6 {
			t.Errorf("scanDigits(s, 0) returns %d, %d, nil", start, end)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte("123456abcdefghijkl")
		//             |<>|
		start, end, err := scanDigits(s, 2)
		if err != nil {
			t.Errorf("scanDigits(s, 2) returns %d, %d, %s", start, end, err)
		}

		if start != 2 || end != 6 {
			t.Errorf("scanDigits(s, 2) returns %d, %d, nil", start, end)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte("123456abcdefghijkl")
		//                |^
		start, end, err := scanDigits(s, 6)
		if err == nil {
			t.Fatalf("scanDigits(s, 6) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 6: expect digit, got 'a'" {
			t.Errorf("scanDigits(s, 6) returns %d, %d, %s", start, end, err)
		}

		if start != 6 || end != 6 {
			t.Errorf("scanDigits(s, 6) returns %d, %d, %s", start, end, err)
		}
	}
}

func TestScanHexDigits(t *testing.T) {
	{
		//           0         1         2
		//           012345678901234567890
		s := []byte("123456abcdefghijkl")
		//           |<-------->|
		start, end, err := scanHexDigits(s, 0)
		if err != nil {
			t.Errorf("scanHexDigits(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 12 {
			t.Errorf("scanHexDigits(s, 0) returns %d, %d, nil", start, end)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte("123456abcdefghijkl")
		//                 |<-->|
		start, end, err := scanHexDigits(s, 6)
		if err != nil {
			t.Errorf("scanHexDigits(s, 6) returns %d, %d, %s", start, end, err)
		}

		if start != 6 || end != 12 {
			t.Errorf("scanHexDigits(s, 6) returns %d, %d, nil", start, end)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte("123456abcdefghijkl")
		//                      |^
		start, end, err := scanHexDigits(s, 12)
		if err == nil {
			t.Fatalf("scanHexDigits(s, 12) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 12: expect hex digit, got 'g'" {
			t.Errorf("scanHexDigits(s, 12) returns %d, %d, %s", start, end, err)
		}

		if start != 12 || end != 12 {
			t.Errorf("scanHexDigits(s, 12) returns %d, %d, %s", start, end, err)
		}
	}
}

func TestScanJsonNumberSuccess(t *testing.T) {
	caseList := scannerCorrectCases{
		"0",
		"-0",

		"12345",
		"-12345",

		"1234.5678",
		"-1234.5678",

		"1234e5678",
		"1234e+5678",
		"1234e-5678",
		"1234E5678",
		"1234E+5678",
		"1234E-5678",
		"-1234e5678",
		"-1234e+5678",
		"-1234e-5678",
		"-1234E5678",
		"-1234E+5678",
		"-1234E-5678",

		"1234.5678e90",
		"1234.5678e+90",
		"1234.5678e-90",
		"1234.5678E90",
		"1234.5678E+90",
		"1234.5678E-90",
		"-1234.5678e90",
		"-1234.5678e+90",
		"-1234.5678e-90",
		"-1234.5678E90",
		"-1234.5678E+90",
		"-1234.5678E-90",
	}

	caseList.On(t, scanJsonNumber)
}

func TestScanJsonNumberFailure(t *testing.T) {
	{
		s := []byte("-")
		start, end, err := scanJsonNumber(s, 0)
		if err == nil {
			t.Fatalf("scanJsonNumber(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 1: expect digit, got 'EOF'" {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 1 {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, nil", start, end)
		}
	}

	{
		s := []byte("-invalid")
		start, end, err := scanJsonNumber(s, 0)
		if err == nil {
			t.Fatalf("scanJsonNumber(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 1: expect digit, got 'i'" {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 1 {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           01234567890
		s := []byte("12345abcde")
		//               |^
		start, end, err := scanJsonNumber(s, 5)
		if err == nil {
			t.Fatalf("scanJsonNumber(s, 5) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 5: expect digit or '-', got 'a'" {
			t.Errorf("scanJsonNumber(s, 5) returns %d, %d, %s", start, end, err)
		}

		if start != 5 || end != 5 {
			t.Errorf("scanJsonNumber(s, 5) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           01234567890
		s := []byte("12345.abcde")
		//          |      ^
		start, end, err := scanJsonNumber(s, 0)
		if err == nil {
			t.Fatalf("scanJsonNumber(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 6: expect digit, got 'a'" {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 6 {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           01234567890
		s := []byte("12345einvalid")
		//          |      ^
		start, end, err := scanJsonNumber(s, 0)
		if err == nil {
			t.Fatalf("scanJsonNumber(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 6: expect digit, got 'i'" {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 6 {
			t.Errorf("scanJsonNumber(s, 0) returns %d, %d, %s", start, end, err)
		}
	}
}

func TestScanJsonStringSuccess(t *testing.T) {
	caseList := scannerCorrectCases{
		`"abc"`,
		`"abc\ndef\fghi\tjkl\rmno\"pqr\\stu"`,
		`"abc\u4daebc"`,
	}

	caseList.On(t, scanJsonString)
}

func TestScanJsonStringFailure(t *testing.T) {
	{
		s := []byte(`invalid`)
		start, end, err := scanJsonString(s, 0)
		if err == nil {
			t.Fatalf("scanJsonString(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 0: expect quote '\"', got 'i'" {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 0 {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           0123456789012345
		s := []byte(`"the quick \uinvalid brown fox"`)
		//                        ^
		start, end, err := scanJsonString(s, 0)
		if err == nil {
			t.Fatalf("scanJsonString(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 13: expect hex digit, got 'i'" {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 13 {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           0123456789012345
		s := []byte(`"the quick \u42invalid brown fox"`)
		//           |            ^~~^
		start, end, err := scanJsonString(s, 0)
		if err == nil {
			t.Fatalf("scanJsonString(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 13: expect 4 hex digits, got '42in'" {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 13 {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           0123456789012345
		s := []byte(`"the quick \invalid brown fox"`)
		//           |           ^
		start, end, err := scanJsonString(s, 0)
		if err == nil {
			t.Fatalf("scanJsonString(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 12: expect escape char, got 'i'" {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 12 {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte(`"the quick brown fox`)
		//           |                   ^
		start, end, err := scanJsonString(s, 0)
		if err == nil {
			t.Fatalf("scanJsonString(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 20: expect quote '\"', got 'EOF'" {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 20 {
			t.Errorf("scanJsonString(s, 0) returns %d, %d, %s", start, end, err)
		}

	}
}

func TestScanJsonArraySuccess(t *testing.T) {
	caseList := scannerCorrectCases{
		"[]",
		"[1, 1,  2, 3,5,8    ,13]",
		"[[[]]]",
		"[[],[],[],[],[]]",
		`[42, 299792458, 3.1415926, -273.15, 8987661788.7,
			6.02214076e23, 6.02214076e+23, 6.62607015e-34,
			"LOREM", "IPSUM", ["LOREM", "IPSUM"], true, null,
			{"gravitation": 6.67430e-11,
			 	"elementary charge": 1.602176634e-19}]`,
	}

	caseList.On(t, scanJsonArray)
}

func TestScanJsonArrayFailure(t *testing.T) {
	{
		s := []byte(`invalid`)
		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 0: expect bracket '[', got 'i'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 0 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		s := []byte("123456")
		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 0: expect bracket '[', got '1'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 0 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		s := []byte("[")
		//          | ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 1: expect value or bracket ']', got 'EOF'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 1 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte("[1, 1,  2, 3,5,8    ,13")
		//          |                       ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 23: expect comma',' or bracket ']', got 'EOF'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 23 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte("[1, 1,  2, 3,5,8    ,13,]")
		//          |                        ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 24: unexpected first char ']'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 24 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte("[1, 1,  2, 3,5,8    ,")
		//          |                     ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 21: array is not close, got 'EOF'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 21 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte("[1, 1,  2, 3,5,8    ,   ")
		//          |                        ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 24: expect value or bracket ']', got 'EOF'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 24 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte("[1, 1,  2, 3,5, invalid ]")
		//          |                ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 16: unexpected first char 'i'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 16 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte("[1, 1,  2, 3,5, 8 invalid ]")
		//          |                  ^

		start, end, err := scanJsonArray(s, 0)
		if err == nil {
			t.Fatalf("scanJsonArray(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 18: expect comma ',' or bracket ']', got 'i'" {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 18 {
			t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
		}
	}
}

func TestScanJsonObjectSuccess(t *testing.T) {
	caseList := scannerCorrectCases{
		`{}`,
		`{"a":1, "b":2, "c":3}`,
		`{"gravitation": 6.67430e-11, "elementary charge": 1.602176634e-19, 
			"fibonacci": [1, 1, 2, 3, 5, 8, 13, 21], "lorem": "ipsum",
			"boolean": true}`,
	}

	caseList.On(t, scanJsonObject)
}

func TestScanVeryDeepArray(t *testing.T) {
	depth := 1000 * 1000
	lBrackets := strings.Repeat("[", depth)
	rBrackets := strings.Repeat("]", depth)
	s := []byte(lBrackets + rBrackets)

	start, end, err := scanJsonArray(s, 0)
	if err != nil {
		t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
	}
}

func TestScanJsonObjectFailure(t *testing.T) {
	{
		s := []byte(`invalid`)
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 0: expect brace '{', got 'i'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 0 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		s := []byte(`123456`)
		//          |^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 0: expect brace '{', got '1'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 0 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		s := []byte(`{`)
		//          | ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 1: expect key string or brace '}', got 'EOF'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 1 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1
		//           01234567890
		s := []byte(`{1: "one"}`)
		//          | ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 1: expect quote '\"', got '1'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 1 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte(`{"one": 1, "two": 2,`)
		//          |                    ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 20: object is not close, got 'EOF'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 20 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte(`{"one": 1, "two": 2,  `)
		//          |                      ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 22: expect key string, got 'EOF'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 22 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           01234567890123456789012345
		s := []byte(`{"one": 1, "two": 2,}`)
		//          |                    ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 20: expect quote '\"', got '}'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 20 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte(`{"one": 1, "two"  `)
		//          |                  ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 18: expect colon ':', got 'EOF'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 18 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2         3
		//           0123456789012345678901234567890
		s := []byte(`{"one": 1, "two"  invalid}`)
		//          |                  ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 18: expect colon ':', got 'i'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 18 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2         3
		//           0123456789012345678901234567890
		s := []byte(`{"one": 1, "two": `)
		//          |                  ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 18: expect value, got 'EOF'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 18 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2         3
		//           0123456789012345678901234567890
		s := []byte(`{"one": 1, "two": invalid}`)
		//          |                  ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 18: unexpected first char 'i'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 18 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2
		//           012345678901234567890
		s := []byte(`{"one": 1, "two": 2`)
		//          |                   ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 19: expect comma ',' or brace '}', got 'EOF'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 19 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}

	{
		//           0         1         2         3
		//           0123456789012345678901234567890
		s := []byte(`{"one": 1, "two": 2 invalid}`)
		//          |                    ^
		start, end, err := scanJsonObject(s, 0)
		if err == nil {
			t.Fatalf("scanJsonObject(s, 0) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 20: expect comma ',' or brace '}', got 'i'" {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 20 {
			t.Errorf("scanJsonObject(s, 0) returns %d, %d, %s", start, end, err)
		}
	}
}

func TestScanVeryDeepObject(t *testing.T) {
	depth := 1000 * 1000
	lBrackets := strings.Repeat(`{"_":`, depth)
	rBrackets := strings.Repeat("}", depth)
	s := []byte(lBrackets + "{}" + rBrackets)

	start, end, err := scanJsonObject(s, 0)
	if err != nil {
		t.Errorf("scanJsonArray(s, 0) returns %d, %d, %s", start, end, err)
	}
}
