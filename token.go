package main

import (
	"strconv"
	"strings"
	"unicode"
)

type tokenType byte

const (
	TOK_EOF    tokenType = 0
	TOK_OPEN   tokenType = 1
	TOK_CLOSE  tokenType = 2
	TOK_DOT    tokenType = 3
	TOK_SYMBOL tokenType = 4
	TOK_FIXNUM tokenType = 5
)

type tokenizer struct {
	input *strings.Reader
	char  rune
	token struct {
		Type  tokenType
		Value string
	}
}

func newTokenizer(input string) *tokenizer {
	self := &tokenizer{input: strings.NewReader(input)}
	self.nextChar()
	self.Next()
	return self
}

func (self *tokenizer) Next() {
	self.token.Value = ""

	// Skip whitespace
	for unicode.IsSpace(self.char) {
		self.nextChar()
	}

	if self.char == '(' {
		self.token.Type = TOK_OPEN
		self.nextChar()

	} else if self.char == ')' {
		self.token.Type = TOK_CLOSE
		self.nextChar()

	} else if self.char == '.' {
		self.token.Type = TOK_DOT
		self.nextChar()

	} else if unicode.IsDigit(self.char) {
		self.token.Type = TOK_FIXNUM

		for {
			self.token.Value += string(self.char)
			self.nextChar()

			if !unicode.IsDigit(self.char) {
				break
			}
		}

	} else if isValidSymbolChar(self.char) {
		self.token.Type = TOK_SYMBOL

		for {
			self.token.Value += string(self.char)
			self.nextChar()

			if !isValidSymbolChar(self.char) {
				break
			}
		}

	} else {
		// End of input
		self.token.Type = TOK_EOF
	}
}

func (self *tokenizer) Type() tokenType {
	return self.token.Type
}

func (self *tokenizer) StringValue() string {
	return self.token.Value
}

func (self *tokenizer) IntValue() int64 {
	var value int64
	var err error

	if value, err = strconv.ParseInt(self.token.Value, 10, 64); err != nil {
		panic(err)
	}

	return value
}

func (self *tokenizer) nextChar() {
	var err error

	if self.char, _, err = self.input.ReadRune(); err != nil {
		self.char = 0
	}
}

func isValidSymbolChar(char rune) bool {
	return !unicode.IsSpace(char) &&
		!unicode.IsControl(char) &&
		!strings.ContainsRune("().", char)
}
