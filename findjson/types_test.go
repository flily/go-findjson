package findjson

import (
	"testing"
)

func TestJsonValueKind(t *testing.T) {
	if f := GetScannerInJNS(JsonValueArray); f == nil {
		t.Errorf("GetScannerOf(JsonValueArray) returns nil")
	}

	if f := GetScannerInJNS(JsonValueAll); f != nil {
		t.Errorf("GetScannerOf(JsonValueAll) returns %v", f)
	}

	j := JsonValueArray
	if j.CanScan(JsonValueArray, NormativeStyle) == nil {
		t.Errorf("CanScan(JsonValueArray) returns nil")
	}

	if f := j.CanScan(JsonValueAll, NormativeStyle); f != nil {
		t.Errorf("CanScan(JsonValueAll) returns %v", f)
	}

	if f := j.CanScan(JsonValueObject, NormativeStyle); f != nil {
		t.Errorf("CanScan(JsonValueObject) returns %v", f)
	}
}

func TestJsonValueKindString(t *testing.T) {
	if v := JsonValueNull.String(); v != "null" {
		t.Errorf("JsonValueNull.String() returns '%s'", v)
	}

	if v := JsonValueBoolean.String(); v != "boolean" {
		t.Errorf("JsonValueBoolean.String() returns '%s'", v)
	}

	if v := JsonValueNumber.String(); v != "number" {
		t.Errorf("JsonValueNumber.String() returns '%s'", v)
	}

	if v := JsonValueString.String(); v != "string" {
		t.Errorf("JsonValueString.String() returns '%s'", v)
	}

	if v := JsonValueArray.String(); v != "array" {
		t.Errorf("JsonValueArray.String() returns '%s'", v)
	}

	if v := JsonValueObject.String(); v != "object" {
		t.Errorf("JsonValueObject.String() returns '%s'", v)
	}

	if v := JsonValueAll.String(); v != "null|boolean|number|string|array|object" {
		t.Errorf("JsonValueAll.String() returns '%s'", v)
	}
}
