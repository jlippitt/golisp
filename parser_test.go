package main

import "testing"

func TestParser(t *testing.T) {
	expected := "(+ . (6 . (5 . ())))"

	actual := parse("(+ 6 5)").String()

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
