package main

import (
	"bytes"
	"testing"
)

func TestOpNil(t *testing.T) {
	code := []byte{
		0x00, // NIL
	}

	returnValue := run(bytes.NewReader(code))

	switch returnValue := returnValue.(type) {
	case *nilCell:
		// Nothing
	default:
		t.Errorf("Expected type Nil, got %s", dump(returnValue))
	}
}

func TestOpLdc(t *testing.T) {
	var intValue int64
	var expectedValue int64 = -4822679049321438

	code := []byte{
		0x01, 0xFF, 0xEE, 0xDD, 0xCC, 0x88, 0x66, 0x44, 0x22, // LDC
	}

	returnValue := run(bytes.NewReader(code))

	switch returnValue := returnValue.(type) {
	case *fixNumCell:
		intValue = returnValue.Value()

		if intValue != expectedValue {
			t.Errorf("Expected %d, got %d", expectedValue, intValue)
		}

	default:
		t.Errorf("Expected type FixNum, got %s", dump(returnValue))
	}
}

func TestOpAdd(t *testing.T) {
	var intValue int64

	code := []byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, // LDC 6
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // LDC 5
		0x02, // ADD
	}

	returnValue := run(bytes.NewReader(code))

	switch returnValue := returnValue.(type) {
	case *fixNumCell:
		intValue = returnValue.Value()

		if intValue != 11 {
			t.Errorf("Expected 11, got %d", intValue)
		}

	default:
		t.Errorf("Expected type FixNum, got %s", dump(returnValue))
	}
}
