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

	expected := []byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // LDC 5
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, // LDC 6
		0x02, // ADD
	}

	actual := generateCode(ast)

	if actual.Len() != len(expected) {
		t.Errorf(
			"Expected code to be %d bytes long, is actually %d bytes long",
			len(expected),
			actual.Len(),
		)
	}

	for i, expectedByte := range expected {
		actualByte, _ := actual.ReadByte()

		if actualByte != expectedByte {
			t.Errorf(
				"Expected %02X at location %d, got %02X",
				expectedByte,
				i,
				actualByte,
			)
		}
	}
}
