package findjson

import (
	"testing"
)

func TestJsonValueKind(t *testing.T) {
	if f := GetScannerOf(JsonValueArray); f == nil {
		t.Errorf("GetScannerOf(JsonValueArray) returns nil")
	}

	if f := GetScannerOf(JsonValueAll); f != nil {
		t.Errorf("GetScannerOf(JsonValueAll) returns %v", f)
	}

	j := JsonValueArray
	if j.CanScan(JsonValueArray) == nil {
		t.Errorf("CanScan(JsonValueArray) returns nil")
	}

	if f := j.CanScan(JsonValueAll); f != nil {
		t.Errorf("CanScan(JsonValueAll) returns %v", j.CanScan(JsonValueAll))
	}

	if f := j.CanScan(JsonValueObject); f != nil {
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
