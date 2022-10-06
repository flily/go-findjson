package findjson

import (
	"strings"
)

// Json value element kind, can be use as a bitmask.
//
// One of the following:
//
//	JsonValueNull
//	JsonValueBoolean
//	JsonValueNumber
//	JsonValueString
//	JsonValueArray
//	JsonValueObject
type JsonValueKind int

const (
	JsonValueNull    = JsonValueKind(1)
	JsonValueBoolean = JsonValueKind(2)
	JsonValueNumber  = JsonValueKind(4)
	JsonValueString  = JsonValueKind(8)
	JsonValueArray   = JsonValueKind(16)
	JsonValueObject  = JsonValueKind(32)
	JsonValueAll     = JsonValueKind(0x00ff) // all values, 0x003f actually
)

// Get JSON value scanner by value kind
type JsonScannerProvider func(JsonValueKind) JsonTokenScanner

func GetScannerInJNS(kind JsonValueKind) JsonTokenScanner {
	switch kind {
	case JsonValueNull:
		return scanJsonLiteral

	case JsonValueBoolean:
		return scanJsonLiteral

	case JsonValueNumber:
		return scanJsonNumber

	case JsonValueString:
		return scanJsonString

	case JsonValueArray:
		return scanJsonArrayJNS

	case JsonValueObject:
		return scanJsonObjectJNS
	}

	return nil
}

func GetScannerInJSS(kind JsonValueKind) JsonTokenScanner {
	switch kind {
	case JsonValueNull:
		return scanJsonLiteral

	case JsonValueBoolean:
		return scanJsonLiteral

	case JsonValueNumber:
		return scanJsonNumber

	case JsonValueString:
		return scanJsonString

	case JsonValueArray:
		return scanJsonArrayJSS

	case JsonValueObject:
		return scanJsonObjectJSS

	}

	return nil
}

func GetScannerOf(kind JsonValueKind, style int) JsonTokenScanner {
	switch style {
	case NormativeStyle:
		return GetScannerInJNS(kind)

	case JavaScriptStyle:
		return GetScannerInJSS(kind)
	}

	return nil
}

func (k JsonValueKind) String() string {
	names := make([]string, 0, 6)
	if k&JsonValueNull != 0 {
		names = append(names, "null")
	}

	if k&JsonValueBoolean != 0 {
		names = append(names, "boolean")
	}

	if k&JsonValueNumber != 0 {
		names = append(names, "number")
	}

	if k&JsonValueString != 0 {
		names = append(names, "string")
	}

	if k&JsonValueArray != 0 {
		names = append(names, "array")
	}

	if k&JsonValueObject != 0 {
		names = append(names, "object")
	}

	return strings.Join(names, "|")
}

func (k JsonValueKind) CanScan(kind JsonValueKind, style int) JsonTokenScanner {
	if k&kind != 0 {
		return GetScannerOf(kind, style)
	}

	return nil
}

func (k JsonValueKind) GetScanner(c byte, style int) JsonTokenScanner {
	switch c {
	case 'n':
		return k.CanScan(JsonValueNull, style)

	case 't', 'f':
		return k.CanScan(JsonValueBoolean, style)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return k.CanScan(JsonValueNumber, style)

	case jsonQuote:
		return k.CanScan(JsonValueString, style)

	case jsonLBracket:
		return k.CanScan(JsonValueArray, style)

	case jsonLBrace:
		return k.CanScan(JsonValueObject, style)

	default:
		return nil
	}
}
