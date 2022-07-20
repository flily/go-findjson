package findjson

import (
	"testing"
)

func TestJsonValueType(t *testing.T) {

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
