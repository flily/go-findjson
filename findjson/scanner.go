package findjson

const (
	JsonBackslash     = '\\'
	JsonSignPositive  = '+'
	JsonSignNegative  = '-'
	JsonExponentUpper = 'E'
	JsonExponentLower = 'e'
	JsonLBracket      = '['
	JsonRBracket      = ']'
	JsonLBrace        = '{'
	JsonRBrace        = '}'
	JsonColon         = ':'
	JsonComma         = ','
	JsonPeriod        = '.'
	JsonQuote         = '"'
	JsonDigitZero     = '0'
	JsonUnicode       = 'u'
)

var (
	JsonLiteralTrue  = []byte("true")
	JsonLiteralFalse = []byte("false")
	JsonLiteralNull  = []byte("null")
)

func bufferStartsWith(s []byte, i int, prefix []byte) bool {
	sl := len(s)
	pl := len(prefix)

	if i+pl > sl {
		return false
	}

	for j := 0; i+j < sl && j < pl; j++ {
		if s[i+j] != prefix[j] {
			return false
		}
	}

	return true
}

func bufferFindSample(s []byte, i int, max int) string {
	l := len(s)
	if i >= l {
		return "EOF"
	}

	if i+max > l {
		return string(s[i:])
	}

	return string(s[i : i+max])
}

func jumpNextNonWhiteSpace(s []byte, i int) int {
	l := len(s)

	for i < l && isWhiteSpace(s[i]) {
		i++
	}

	return i
}

/*
 * scanX functions
 * scan(s []byte, i int) (start int, next int, err error)
 * scan JSON grammar element.
 * return start of this element, next position after this element, and error.
 */

type JsonTokenScanner func(s []byte, i int) (tokenStart int, tokenNext int, err error)

func scanJsonLiteral(s []byte, i int) (int, int, error) {
	var err error
	j := i

	if bufferStartsWith(s, i, JsonLiteralNull) {
		j += 4

	} else if bufferStartsWith(s, i, JsonLiteralTrue) {
		j += 4

	} else if bufferStartsWith(s, i, JsonLiteralFalse) {
		j += 5

	} else {
		v := bufferFindSample(s, i, 5)
		err = NewJsonError(i, "expect null, true or false, got '%s'", v)
	}

	return i, j, err
}

func tryScanCharSet(s []byte, i int, check CharSetChecker) int {
	l, j := len(s), i

	for j < l && check(s[j]) {
		j++
	}

	return j
}

func scanDigits(s []byte, i int) (int, int, error) {
	var err error
	j := tryScanCharSet(s, i, isDigit)

	if j <= i {
		v := bufferFindSample(s, i, 1)
		err = NewJsonError(i, "expect digit, got '%s'", v)
	}

	return i, j, err
}

func scanHexDigits(s []byte, i int) (int, int, error) {
	var err error
	j := tryScanCharSet(s, i, isHexDigit)

	if j <= i {
		v := bufferFindSample(s, i, 5)
		err = NewJsonError(i, "expect hex digit, got '%s'", v)
	}

	return i, j, err
}

func scanJsonNumber(s []byte, i int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	found_neg_sign := false

	if s[j] == JsonSignNegative {
		j++
		found_neg_sign = true
	}

	if j >= l {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect digit, got '%s'", v)
		return i, j, err
	}

	if s[j] == JsonDigitZero {
		// ZERO or floats
		// JSON do not support integers in octal form (e.g. 0777), nor hexadecimal (e.g. 0xdead)
		// JSON do not support float number without leading digit (e.g. .233)
		j++

	} else if isNonZeroDigit(s[j]) {
		_, j, err = scanDigits(s, j)

	} else {
		v := bufferFindSample(s, j, 1)

		if found_neg_sign {
			err = NewJsonError(j, "expect digit, got '%s'", v)

		} else {
			// May not reach here if called via FIRST_SET tables.
			err = NewJsonError(j, "expect digit or '-', got '%s'", v)
		}

		return i, j, err
	}

	if j < l && s[j] == JsonPeriod {
		// fraction
		_, j, err = scanDigits(s, j+1)
		if err != nil {
			return i, j, err
		}
	}

	if j < l && (s[j] == JsonExponentUpper || s[j] == JsonExponentLower) {
		// exponent
		j++
		if j < l && (s[j] == JsonSignPositive || s[j] == JsonSignNegative) {
			j++
		}

		_, j, err = scanDigits(s, j)
		if err != nil {
			return i, j, err
		}
	}

	return i, j, err
}

func scanJsonString(s []byte, i int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	quote_close := false

	if s[j] != JsonQuote {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect quote '\"', got '%s'", v)
		return i, j, err
	}

	j++ // skip quote
	for j < l {
		c1 := s[j]
		j++

		if c1 == JsonBackslash {
			c2 := s[j]
			j++

			if isEscapeChar(c2) {
				// escape char

			} else if c2 == JsonUnicode {
				var nj int
				_, nj, err = scanHexDigits(s, j)
				if err != nil {
					break
				}

				if nj-j >= 4 {
					j += 4

				} else {
					v := bufferFindSample(s, j, 4)
					err = NewJsonError(j, "expect 4 hex digits, got '%s'", v)
					break
				}

			} else {
				v := bufferFindSample(s, j, 1)
				err = NewJsonError(j, "expect escape char, got '%s'", v)
				break
			}

		} else if c1 == JsonQuote {
			quote_close = true
			break
		}
	}

	if err == nil && !quote_close {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect quote '\"', got '%s'", v)
	}

	return i, j, err
}

func scanJsonArray(s []byte, i int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	bracket_close := false

	if s[j] != JsonLBracket {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect '[', got '%s'", v)
		return i, j, err
	}

	j = jumpNextNonWhiteSpace(s, j+1)
	if s[j] == JsonRBracket {
		return i, j + 1, nil
	}

	for j < l {
		j = jumpNextNonWhiteSpace(s, j)
		c := s[j]

		scanner := scannerByFirstSet(c)
		if scanner == nil {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "unexpected char '%s'", v)
			break
		}

		_, j, err = scanner(s, j)
		if err != nil {
			break
		}

		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect ',' or ']', got '%s'", v)

		} else if s[j] == JsonComma {
			j += 1
			continue

		} else if s[j] == JsonRBracket {
			j += 1
			bracket_close = true

		} else {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "unexpected char '%s'", v)
		}

		break
	}

	if err == nil && !bracket_close {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "array is not close, got '%s'", v)
	}

	return i, j, err
}

func scanJsonObject(s []byte, i int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	brace_close := false

	if s[j] != JsonLBrace {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect '{', got '%s'", v)
		return i, j, err
	}

	j++
	for j < l {
		j = jumpNextNonWhiteSpace(s, j)
		if s[j] == JsonRBrace {
			j++
			brace_close = true
			break
		}

		_, j, err = scanJsonString(s, j)
		if err != nil {
			return i, j, err
		}

		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect ':', got '%s'", v)
			break

		} else if s[j] != JsonColon {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect ':', got '%s'", v)
			break
		}

		j = jumpNextNonWhiteSpace(s, j+1)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect value, got '%s'", v)
			break
		}

		c := s[j]
		scanner := scannerByFirstSet(c)
		if scanner == nil {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "unexpected char '%s'", v)
			break
		}

		_, j, err = scanner(s, j)
		if err != nil {
			break
		}

		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect ',' or '}', got '%s'", v)

		} else if s[j] == JsonComma {
			j += 1
			continue

		} else if s[j] == JsonRBrace {
			j += 1
			brace_close = true

		} else {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "unexpected char '%s'", v)
		}

		break
	}

	if err == nil && !brace_close {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "object is not close, got '%s'", v)
	}

	return i, j, err
}

func scannerByFirstSet(c byte) JsonTokenScanner {
	switch c {
	case 'n', 't', 'f':
		return scanJsonLiteral

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return scanJsonNumber

	case JsonQuote:
		return scanJsonString

	case JsonLBracket:
		return scanJsonArray

	case JsonLBrace:
		return scanJsonObject

	default:
		return nil
	}
}
