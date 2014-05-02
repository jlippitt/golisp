package main

import "testing"

func checkResult(t *testing.T, code string, expected cell) {
	var actual cell
	var err error

	actual, err = Execute(code)

	if err != nil {
		t.Errorf("Error occured while executing code: %s\n", err.Error())
	}

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}
}

func TestOperators(t *testing.T) {
	checkResult(t, "(+ 6 5)", newFixNumCell(11))
}

func TestIfCondition(t *testing.T) {
	checkResult(t, "(if (= 10 10) 1 2)", newFixNumCell(1))
	checkResult(t, "(if (= 10 11) 1 2)", newFixNumCell(2))
	checkResult(t, "(if (= 10 10) (* 2 2) 2)", newFixNumCell(4))
	checkResult(t, "(if (= 10 11) 1 (* 4 4))", newFixNumCell(16))
}
