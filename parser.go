package main

import "fmt"

func parse(input string) cell {
	var parseExpression, parseCons func() cell

	tokenizer := newTokenizer(input)

	parseExpression = func() cell {
		var value cell

		switch tokenizer.Type() {
		case TOK_OPEN:
			value = parseCons()
		case TOK_SYMBOL:
			switch symbol := tokenizer.StringValue(); symbol {
			case "T":
				value = newTrueCell()
			case "nil":
				value = newNilCell()
			default:
				value = newSymbolCell(symbol)
			}
		case TOK_FIXNUM:
			value = newFixNumCell(tokenizer.IntValue())
		default:
			panic(fmt.Sprintf("Unexpected token: %d", tokenizer.Type()))
		}

		return value
	}

	parseCons = func() cell {
		var cons cell = newNilCell()
		var it *cell = &cons

		for {
			tokenizer.Next()

			if tokenizer.Type() != TOK_CLOSE {
				pushBack(&it, parseExpression())
			} else {
				break
			}
		}

		return cons
	}

	var ast cell = newNilCell()
	var it *cell = &ast

	for tokenizer.Type() != TOK_EOF {
		pushBack(&it, parseExpression())
		tokenizer.Next()
	}

	return ast
}
