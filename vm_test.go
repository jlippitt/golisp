package main

import (
	"testing"
)

func TestOpNil(t *testing.T) {
	code := newConsCell(
		newOpCell(OP_NIL, nil),
		newConsCell(
			newOpCell(OP_HALT, nil),
			newNilCell(),
		),
	)

	expected := newNilCell()

	actual := run(code)

	if dump(actual) != dump(expected) {
		t.Errorf("Expected %s but got %s", dump(expected), dump(actual))
	}
}

func TestOpLdc(t *testing.T) {
	code := newConsCell(
		newOpCell(OP_LDC, newFixNumCell(-52)),
		newConsCell(
			newOpCell(OP_HALT, nil),
			newNilCell(),
		),
	)

	expected := newFixNumCell(-52)

	actual := run(code)

	if dump(actual) != dump(expected) {
		t.Errorf("Expected %s but got %s", dump(expected), dump(actual))
	}
}

func TestOpAdd(t *testing.T) {
	code := newConsCell(
		newOpCell(OP_LDC, newFixNumCell(5)),
		newConsCell(
			newOpCell(OP_LDC, newFixNumCell(6)),
			newConsCell(
				newOpCell(OP_ADD, nil),
				newConsCell(
					newOpCell(OP_HALT, nil),
					newNilCell(),
				),
			),
		),
	)

	expected := newFixNumCell(11)

	actual := run(code)

	if dump(actual) != dump(expected) {
		t.Errorf("Expected %s but got %s", dump(expected), dump(actual))
	}
}
