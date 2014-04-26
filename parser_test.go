package main

import "testing"

func TestParser(t *testing.T) {
	expected := "(+ . (6 . (5 . ())))"

	actual := dump(parse("(+ 6 5)"))

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
