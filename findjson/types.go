package findjson

import "strings"

// Json value element kind, can be use as a bitmask.
//
// One of the following:
//  JsonValueNull
//  JsonValueBoolean
//  JsonValueNumber
//  JsonValueString
//  JsonValueArray
//  JsonValueObject
type JsonValueKind int

func GetScannerOf(k JsonValueKind) JsonTokenScanner {
	switch k {
	case JsonValueNull:
		return scanJsonLiteral

	case JsonValueBoolean:
		return scanJsonLiteral

	case JsonValueNumber:
		return scanJsonNumber

	case JsonValueString:
		return scanJsonString

	case JsonValueArray:
		return scanJsonArray

	case JsonValueObject:
		return scanJsonObject
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

func (k JsonValueKind) CanScan(kind JsonValueKind) JsonTokenScanner {
	if k&kind != 0 {
		return GetScannerOf(kind)
	}

	return nil
}

func (k JsonValueKind) GetScanner(c byte) JsonTokenScanner {
	switch c {
	case 'n':
		return k.CanScan(JsonValueNull)

	case 't', 'f':
		return k.CanScan(JsonValueBoolean)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return k.CanScan(JsonValueNumber)

	case jsonQuote:
		return k.CanScan(JsonValueString)

	case jsonLBracket:
		return k.CanScan(JsonValueArray)

	case jsonLBrace:
		return k.CanScan(JsonValueObject)

	default:
		return nil
	}
}

var (
	JsonValueNull    = JsonValueKind(1)
	JsonValueBoolean = JsonValueKind(2)
	JsonValueNumber  = JsonValueKind(4)
	JsonValueString  = JsonValueKind(8)
	JsonValueArray   = JsonValueKind(16)
	JsonValueObject  = JsonValueKind(32)
	JsonValueAll     = JsonValueKind(0x00ff) // all values, 0x003f actually
)
