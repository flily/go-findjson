package findjson

const (
	jsonCharsetWhiteSpace    = 0x01
	jsonCharsetEscapeChars   = 0x02
	jsonCharsetDigits        = 0x04
	jsonCharsetDigitsNonZero = 0x08
	jsonCharsetHexDigits     = 0x10
)

const (
	// abbreviations
	jWSP = jsonCharsetWhiteSpace
	jDGT = jsonCharsetDigits | jsonCharsetDigitsNonZero | jsonCharsetHexDigits
	jHEX = jsonCharsetHexDigits
	jESC = jsonCharsetEscapeChars
	jCH0 = jsonCharsetDigits | jsonCharsetHexDigits
	jCHb = jsonCharsetHexDigits | jsonCharsetEscapeChars
	jCHf = jsonCharsetHexDigits | jsonCharsetEscapeChars
)

var charmap = [256]byte{
	// 0     1     2     3     4     5     6     7
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x00
	// 8     9     A     B     C     D     E     F
	0x00, jWSP, jWSP, 0x00, 0x00, jWSP, 0x00, 0x00, // 0x08
	// 0	 1     2     3     4     5     6     7
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x10
	// 8     9     A     B     C     D     E     F
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x18
	// 0     1     2     3     4     5     6     7
	jWSP, 0x00, jESC, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x20
	// 8     9     A     B     C     D     E     F
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, jESC, // 0x28
	// 0     1     2     3     4     5     6     7
	jCH0, jDGT, jDGT, jDGT, jDGT, jDGT, jDGT, jDGT, // 0x30
	// 8     9     A     B     C     D     E     F
	jDGT, jDGT, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x38
	// 0     1     2     3     4     5     6     7
	0x00, jHEX, jHEX, jHEX, jHEX, jHEX, jHEX, 0x00, // 0x40
	// 8     9     A     B     C     D     E     F
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x48
	// 0     1     2     3     4     5     6     7
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x50
	// 8     9     A     B     C     D     E     F
	0x00, 0x00, 0x00, 0x00, jESC, 0x00, 0x00, 0x00, // 0x58
	// 0     1     2     3     4     5     6     7
	0x00, jHEX, jCHb, jHEX, jHEX, jHEX, jCHf, 0x00, // 0x60
	// 8     9     A     B     C     D     E     F
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, jESC, 0x00, // 0x68
	// 0     1     2     3     4     5     6     7
	0x00, 0x00, jESC, 0x00, jESC, 0x00, 0x00, 0x00, // 0x70
	// 8     9     A     B     C     D     E     F
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x78
	// 0x80 - 0xFF
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x80
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x88
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x90
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0x98
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xA0
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xA8
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xB0
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xB8
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xC0
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xC8
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xD0
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xD8
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xE0
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xE8
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xF0
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 0xF8
}

type charSetChecker func(byte) bool

func isWhiteSpace(c byte) bool {
	return charmap[c]&jsonCharsetWhiteSpace != 0
}

func isEscapeChar(c byte) bool {
	return charmap[c]&jsonCharsetEscapeChars != 0
}

func isDigit(c byte) bool {
	return charmap[c]&jsonCharsetDigits != 0
}

func isHexDigit(c byte) bool {
	return charmap[c]&jsonCharsetHexDigits != 0
}

func isNonZeroDigit(c byte) bool {
	return charmap[c]&jsonCharsetDigitsNonZero != 0
}
