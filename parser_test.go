package main

import "testing"

func TestParser(t *testing.T) {
	expected := newConsCell(
		newSymbolCell("+"),
		newConsCell(
			newFixNumCell(6),
			newConsCell(
				newFixNumCell(5),
				newNilCell(),
			),
		),
	)

	actual := parse("(+ 6 5)")

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
