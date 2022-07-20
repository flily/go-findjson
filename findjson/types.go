package findjson

import "strings"

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
		names = append(names, "JsonValueNull")
	}

	if k&JsonValueBoolean != 0 {
		names = append(names, "JsonValueBoolean")
	}

	if k&JsonValueNumber != 0 {
		names = append(names, "JsonValueNumber")
	}

	if k&JsonValueString != 0 {
		names = append(names, "JsonValueString")
	}

	if k&JsonValueArray != 0 {
		names = append(names, "JsonValueArray")
	}

	if k&JsonValueObject != 0 {
		names = append(names, "JsonValueObject")
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

	case JsonQuote:
		return k.CanScan(JsonValueString)

	case JsonLBracket:
		return k.CanScan(JsonValueArray)

	case JsonLBrace:
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
