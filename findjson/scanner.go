package findjson

const (
	jsonBackslash     = '\\'
	jsonSignPositive  = '+'
	jsonSignNegative  = '-'
	jsonExponentUpper = 'E'
	jsonExponentLower = 'e'
	jsonLBracket      = '['
	jsonRBracket      = ']'
	jsonLBrace        = '{'
	jsonRBrace        = '}'
	jsonColon         = ':'
	jsonComma         = ','
	jsonPeriod        = '.'
	jsonQuote         = '"'
	jsonDigitZero     = '0'
	jsonUnicode       = 'u'
)

var (
	jsonLiteralTrue  = []byte("true")
	jsonLiteralFalse = []byte("false")
	jsonLiteralNull  = []byte("null")
)

func bufferStartsWith(s []byte, i int, prefix []byte) (bool, int) {
	sl := len(s)
	pl := len(prefix)
	j := 0

	for j = 0; i+j < sl && j < pl; j++ {
		if s[i+j] != prefix[j] {
			return false, j
		}
	}

	return j >= pl, j
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

	if i >= l {
		return l
	}

	for i < l && isWhiteSpace(s[i]) {
		i++
	}

	return i
}

// scanX functions
//
// Scan JSON grammar element.
// If success, returns start of this element, position just after the end of this element,
// and a nil error.
//
// If buffer partially match the pattern, tokenNext will be just after the position that matched.
//
// If buffer do not matched even in the first set, tokenNext will be equal to tokenStart.
type JsonTokenScanner func(s []byte, i int) (tokenStart int, tokenNext int, err error)

func scanJsonLiteral(s []byte, i int) (int, int, error) {
	var err error
	var m bool // matched
	var l int  // length
	j := i

	if m, l = bufferStartsWith(s, i, jsonLiteralNull); l > 0 {
		j += l

	} else if m, l = bufferStartsWith(s, i, jsonLiteralTrue); l > 0 {
		j += l

	} else if m, l = bufferStartsWith(s, i, jsonLiteralFalse); l > 0 {
		j += l

	}

	if !m {
		v := bufferFindSample(s, i, 5)
		err = NewJsonError(i, "expect null, true or false, got '%s'", v)
	}

	return i, j, err
}

func tryScanCharSet(s []byte, i int, check charSetChecker) int {
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
		v := bufferFindSample(s, i, 1)
		err = NewJsonError(i, "expect hex digit, got '%s'", v)
	}

	return i, j, err
}

func scanJsonNumber(s []byte, i int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	foundNegativeSign := false

	if s[j] == jsonSignNegative {
		j++
		foundNegativeSign = true
	}

	if j >= l {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect digit, got '%s'", v)
		return i, j, err
	}

	if s[j] == jsonDigitZero {
		// ZERO or floats
		// JSON do not support integers in octal form (e.g. 0777), nor hexadecimal (e.g. 0xdead)
		// JSON do not support float number without leading digit (e.g. .233)
		j++

	} else if isNonZeroDigit(s[j]) {
		_, j, err = scanDigits(s, j)

	} else {
		v := bufferFindSample(s, j, 1)

		if foundNegativeSign {
			err = NewJsonError(j, "expect digit, got '%s'", v)

		} else {
			// May not reach here if called via FIRST SET tables.
			err = NewJsonError(j, "expect digit or '-', got '%s'", v)
		}

		return i, j, err
	}

	if j < l && s[j] == jsonPeriod {
		// fraction
		_, j, err = scanDigits(s, j+1)
		if err != nil {
			return i, j, err
		}
	}

	if j < l && (s[j] == jsonExponentUpper || s[j] == jsonExponentLower) {
		// exponent
		j++
		if j < l && (s[j] == jsonSignPositive || s[j] == jsonSignNegative) {
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
	quoteClosed := false

	if s[j] != jsonQuote {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect quote '\"', got '%s'", v)
		return i, j, err
	}

	j++ // skip quote
	for j < l {
		c1 := s[j]
		j++

		if c1 == jsonBackslash {
			c2 := s[j]

			if isEscapeChar(c2) {
				j++
				// escape char

			} else if c2 == jsonUnicode {
				var nj int
				j++
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

		} else if c1 == jsonQuote {
			quoteClosed = true
			break
		}
	}

	if err == nil && !quoteClosed {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect quote '\"', got '%s'", v)
	}

	return i, j, err
}

func scanJsonArray(s []byte, i int, style int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	bracketClosed := false

	if s[j] != jsonLBracket {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect bracket '[', got '%s'", v)
		return i, j, err
	}

	j = jumpNextNonWhiteSpace(s, j+1)
	if j >= l {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect value or bracket ']', got '%s'", v)
		return i, j, err

	} else if s[j] == jsonRBracket {
		return i, j + 1, nil
	}

	for j < l {
		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err := NewJsonError(j, "expect value or bracket ']', got '%s'", v)
			return i, j, err

		} else if s[j] == jsonRBracket && style == JavaScriptStyle {
			j += 1
			bracketClosed = true
			break
		}

		_, j, err = scanJsonValueByFirstSet(s, j, JsonValueAll, style)
		if err != nil {
			break
		}

		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect comma',' or bracket ']', got '%s'", v)

		} else if s[j] == jsonComma {
			j += 1
			continue

		} else if s[j] == jsonRBracket {
			j += 1
			bracketClosed = true

		} else {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect comma ',' or bracket ']', got '%s'", v)
		}

		break
	}

	if err == nil && !bracketClosed {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "array is not close, got '%s'", v)
	}

	return i, j, err
}

// Scan JSON array in JSON style, the trailing comma is NOT ALLOWED.
func scanJsonArrayJNS(s []byte, i int) (int, int, error) {
	return scanJsonArray(s, i, NormativeStyle)
}

// Scan JSON array in JavaScript style, the trailing comma is ALLOWED.
func scanJsonArrayJSS(s []byte, i int) (int, int, error) {
	return scanJsonArray(s, i, JavaScriptStyle)
}

func scanJsonObject(s []byte, i int, style int) (int, int, error) {
	var err error
	l := len(s)
	j := i
	braceClosed := false

	if s[j] != jsonLBrace {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect brace '{', got '%s'", v)
		return i, j, err
	}

	j = jumpNextNonWhiteSpace(s, j+1)
	if j >= l {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "expect key string or brace '}', got '%s'", v)
		return i, j, err

	} else if s[j] == jsonRBrace {
		return i, j + 1, nil
	}

	for j < l {
		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect key string, got '%s'", v)
			return i, j, err

		} else if s[j] == jsonRBrace && style == JavaScriptStyle {
			j += 1
			braceClosed = true
			break
		}

		_, j, err = scanJsonString(s, j)
		if err != nil {
			return i, j, err
		}

		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect colon ':', got '%s'", v)
			break

		} else if s[j] != jsonColon {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect colon ':', got '%s'", v)
			break
		}

		j = jumpNextNonWhiteSpace(s, j+1)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect value, got '%s'", v)
			break
		}

		_, j, err = scanJsonValueByFirstSet(s, j, JsonValueAll, style)
		if err != nil {
			break
		}

		j = jumpNextNonWhiteSpace(s, j)
		if j >= l {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect comma ',' or brace '}', got '%s'", v)

		} else if s[j] == jsonComma {
			j += 1
			continue

		} else if s[j] == jsonRBrace {
			j += 1
			braceClosed = true

		} else {
			v := bufferFindSample(s, j, 1)
			err = NewJsonError(j, "expect comma ',' or brace '}', got '%s'", v)
		}

		break
	}

	if err == nil && !braceClosed {
		v := bufferFindSample(s, j, 1)
		err = NewJsonError(j, "object is not close, got '%s'", v)
	}

	return i, j, err
}

// Scan JSON Object in JSON style, the trailing comma is NOT ALLOWED.
func scanJsonObjectJNS(s []byte, i int) (int, int, error) {
	return scanJsonObject(s, i, NormativeStyle)
}

// Scan JSON Object in JavaScript style, the trailing comma is ALLOWED.
func scanJsonObjectJSS(s []byte, i int) (int, int, error) {
	return scanJsonObject(s, i, JavaScriptStyle)
}

func scanJsonValueByFirstSet(s []byte, i int, kind JsonValueKind, style int) (int, int, error) {
	c := s[i]
	scanner := kind.GetScanner(c, style)
	if scanner == nil {
		v := bufferFindSample(s, i, 1)
		err := NewJsonError(i, "unexpected first char '%s'", v)
		return i, i, err
	}

	return scanner(s, i)
}
