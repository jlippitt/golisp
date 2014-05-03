package main

import (
	"testing"
)

func TestOpNil(t *testing.T) {
	code := newConsCell(
		newOpCell(opNil, nil),
		newConsCell(
			newOpCell(opStop, nil),
			newNilCell(),
		),
	)

	expected := newNilCell()

	actual := run(code)

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestOpLdc(t *testing.T) {
	code := newConsCell(
		newOpCell(opLdc, newFixNumCell(-52)),
		newConsCell(
			newOpCell(opStop, nil),
			newNilCell(),
		),
	)

	expected := newFixNumCell(-52)

	actual := run(code)

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestOpAdd(t *testing.T) {
	code := newConsCell(
		newOpCell(opLdc, newFixNumCell(5)),
		newConsCell(
			newOpCell(opLdc, newFixNumCell(6)),
			newConsCell(
				newOpCell(opAdd, nil),
				newConsCell(
					newOpCell(opStop, nil),
					newNilCell(),
				),
			),
		),
	)

	expected := newFixNumCell(11)

	actual := run(code)

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
