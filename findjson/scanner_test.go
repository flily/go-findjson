package findjson

import (
	"fmt"
	"testing"
)

func TestBufferStartsWith(t *testing.T) {
	//                0        10        20        30        40
	buffer := []byte("the quick brown fox jumps over the lazy dog")

	if !bufferStartsWith(buffer, 0, []byte("the")) {
		t.Errorf("bufferStartsWith(buffer, 0, \"the\") returns false")
	}

	if bufferStartsWith(buffer, 0, []byte("THE")) {
		t.Errorf("bufferStartsWith(buffer, 0, \"THE\") returns true")
	}

	if !bufferStartsWith(buffer, 4, []byte("quick")) {
		t.Errorf("bufferStartsWith(buffer, 4, \"quick\") returns false")
	}

	if bufferStartsWith(buffer, 4, []byte("QUICK")) {
		t.Errorf("bufferStartsWith(buffer, 4, \"QUICK\") returns true")
	}

	if bufferStartsWith(buffer, 10, []byte("brother")) {
		t.Errorf("bufferStartsWith(buffer, 10, \"brother\") returns true")
	}

	if !bufferStartsWith(buffer, 40, []byte("dog")) {
		t.Errorf("bufferStartsWith(buffer, 40, \"dog\") returns false")
	}

	if bufferStartsWith(buffer, 40, []byte("doggie")) {
		t.Errorf("bufferStartsWith(buffer, 40, \"doggie\") returns true")
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
		s := []byte(` true false null tr`)
		start, end, err := scanJsonLiteral(s, 17)
		if err == nil {
			t.Fatalf("scanJsonLiteral(s, 17) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 17: expect null, true or false, got 'tr'" {
			t.Errorf("scanJsonLiteral(s, 17) returns %d, %d, %s", start, end, err)
		}

		if start != 17 || end != 17 {
			t.Errorf("scanJsonLiteral(s, 17) returns %d, %d, nil", start, end)
		}
	}

	{
		s := []byte(` true false null tr`)
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
		s := []byte(` true false null tr`)
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
	s := []byte("12345abcdefghij")

	{
		start, end, err := scanDigits(s, 0)
		if err != nil {
			t.Errorf("scanDigits(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 5 {
			t.Errorf("scanDigits(s, 0) returns %d, %d, nil", start, end)
		}
	}

	{
		start, end, err := scanDigits(s, 2)
		if err != nil {
			t.Errorf("scanDigits(s, 2) returns %d, %d, %s", start, end, err)
		}

		if start != 2 || end != 5 {
			t.Errorf("scanDigits(s, 2) returns %d, %d, nil", start, end)
		}
	}

	{
		start, end, err := scanDigits(s, 6)
		if err == nil {
			t.Fatalf("scanDigits(s, 6) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 6: expect digit, got 'b'" {
			t.Errorf("scanDigits(s, 6) returns %d, %d, %s", start, end, err)
		}

		if start != 6 || end != 6 {
			t.Errorf("scanDigits(s, 6) returns %d, %d, %s", start, end, err)
		}
	}
}

func TestScanHexDigits(t *testing.T) {
	s := []byte("123456abcdefghijkl")

	{
		start, end, err := scanHexDigits(s, 0)
		if err != nil {
			t.Errorf("scanHexDigits(s, 0) returns %d, %d, %s", start, end, err)
		}

		if start != 0 || end != 12 {
			t.Errorf("scanHexDigits(s, 0) returns %d, %d, nil", start, end)
		}
	}

	{
		start, end, err := scanHexDigits(s, 6)
		if err != nil {
			t.Errorf("scanHexDigits(s, 6) returns %d, %d, %s", start, end, err)
		}

		if start != 6 || end != 12 {
			t.Errorf("scanHexDigits(s, 6) returns %d, %d, nil", start, end)
		}
	}

	{
		start, end, err := scanHexDigits(s, 12)
		if err == nil {
			t.Fatalf("scanHexDigits(s, 12) returns %d, %d, nil", start, end)
		}

		if err.Error() != "JSON error at 12: expect hex digit, got 'ghijk'" {
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
		s := []byte("12345abcde")
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
		s := []byte("12345.abcde")
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
		s := []byte("12345einvalid")
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

func TestScanJsonString(t *testing.T) {
	caseList := scannerCorrectCases{
		`"abc"`,
		`"abc\ndef\fghi\tjkl\rmno\"pqr\\stu"`,
		`"abc\u4daebc"`,
	}

	caseList.On(t, scanJsonString)
}

func TestScanJsonArray(t *testing.T) {
	caseList := scannerCorrectCases{
		"[]",
		"[1, 1,  2, 3,5,8    ,13]",
		"[[[]]]",
		"[[],[],[],[],[]]",
		`[42, 299792458, 3.1415926, -273.15, 8987661788.7,
			6.02214076e23, 6.02214076e+23, 6.62607015e-34,
			"LOREM", "IPSUM", ["LOREM", "IPSUM"], true,
			{"gravitation": 6.67430e-11,
			 	"elementary charge": 1.602176634e-19}]`,
	}

	caseList.On(t, scanJsonArray)
}

func TestScanJsonObject(t *testing.T) {
	caseList := scannerCorrectCases{
		`{}`,
		`{"a":1, "b":2, "c":3}`,
		`{"gravitation": 6.67430e-11, "elementary charge": 1.602176634e-19, 
			"fibonacci": [1, 1, 2, 3, 5, 8, 13, 21], "lorem": "ipsum",
			"boolean": true}`,
	}

	caseList.On(t, scanJsonObject)
}
