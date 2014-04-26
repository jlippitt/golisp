package main

import (
	"bytes"
	"testing"
)

func TestAddition(t *testing.T) {
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
