package main

import "testing"

func TestParser(t *testing.T) {
	expected := newConsCell(
		newConsCell(
			newSymbolCell("-"),
			newConsCell(
				newFixNumCell(3),
				newConsCell(
					newFixNumCell(2),
					newNilCell(),
				),
			),
		),
		newConsCell(
			newConsCell(
				newSymbolCell("+"),
				newConsCell(
					newFixNumCell(6),
					newConsCell(
						newFixNumCell(5),
						newNilCell(),
					),
				),
			),
			newNilCell(),
		),
	)

	actual := parse("(- 3 2)\n(+ 6 5)")

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
