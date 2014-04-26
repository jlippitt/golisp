package main

import "testing"

func TestTokenizer(t *testing.T) {
	tokenizer := newTokenizer("(+ 6 5)")

	expect := func(tokenType tokenType, value string) {
		if tokenizer.Type() != tokenType {
			t.Errorf("Expected type %d but got %d", tokenType, tokenizer.Type())
		}

		if tokenizer.StringValue() != value {
			t.Errorf("Expected value %s but got %s", value, tokenizer.StringValue())
		}
	}

	expect(TOK_OPEN, "")
	tokenizer.Next()
	expect(TOK_SYMBOL, "+")
	tokenizer.Next()
	expect(TOK_FIXNUM, "6")
	tokenizer.Next()
	expect(TOK_FIXNUM, "5")
	tokenizer.Next()
	expect(TOK_CLOSE, "")
}
