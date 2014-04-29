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

	st := newSymbolTable()

	expected := newConsCell(
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

	actual := generateCode(ast, st)

	if dump(actual) != dump(expected) {
		t.Errorf("Expected %s but got %s", dump(expected), dump(actual))
	}
}
