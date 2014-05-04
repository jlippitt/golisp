package main

import "testing"

func TestCodeGeneration(t *testing.T) {
	ast := newConsCell(
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

	expected := newConsCell(
		newOpCell(opLdc, newFixNumCell(2)),
		newConsCell(
			newOpCell(opLdc, newFixNumCell(3)),
			newConsCell(
				newOpCell(opSub, nil),
				newConsCell(
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
				),
			),
		),
	)

	actual := generateCode(ast)

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
