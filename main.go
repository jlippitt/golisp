package main

import "log"

func main() {
	code := "(+ 6 5)"

	tokenizer := newTokenizer(code)

	for tokenizer.Type() != TOK_EOF {
		log.Printf("%02X %s\n", tokenizer.Type(), tokenizer.StringValue())
		tokenizer.Next()
	}
}
