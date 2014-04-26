package main

import "testing"

func TestCodeGeneration(t *testing.T) {
	ast := newConsCell(
		newSymbolCell("+"),
		newConsCell(
			newFixNumCell(6),
			newConsCell(
				newFixNumCell(5),
				newNilCell(),
			),
		),
	)

	expected := newConsCell(
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

	actual := generateCode(ast)

	if dump(actual) != dump(expected) {
		t.Errorf("Expected %s but got %s", dump(expected), dump(actual))
	}
}
