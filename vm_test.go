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

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
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

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestOpAdd(t *testing.T) {
	st := newSymbolTable()

	code := newConsCell(
		newOpCell(OP_NIL, nil),
		newConsCell(
			newOpCell(OP_LDC, newFixNumCell(5)),
			newConsCell(
				newOpCell(OP_CONS, nil),
				newConsCell(
					newOpCell(OP_LDC, newFixNumCell(6)),
					newConsCell(
						newOpCell(OP_CONS, nil),
						newConsCell(
							newOpCell(OP_LDC, st.Get("+")),
							newConsCell(
								newOpCell(OP_AP, nil),
								newConsCell(
									newOpCell(OP_HALT, nil),
									newNilCell(),
								),
							),
						),
					),
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
