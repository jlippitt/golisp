package main

func parse(input string) cell {
	return parseExpression(newTokenizer(input))
}

func parseExpression(tokenizer *tokenizer) cell {
	var value cell

	switch tokenizer.Type() {
	case TOK_OPEN:
		value = parseCons(tokenizer)
	case TOK_SYMBOL:
		value = newSymbolCell(tokenizer.StringValue())
	case TOK_FIXNUM:
		value = newFixNumCell(tokenizer.IntValue())
	default:
		panic("Unexpected token")
	}

	return value
}

func parseCons(tokenizer *tokenizer) cell {
	var cons cell = newNilCell()
	var current *cell = &cons

	for {
		tokenizer.Next()

		if tokenizer.Type() != TOK_CLOSE {
			var newValue cell = newConsCell(parseExpression(tokenizer), newNilCell())

			switch cons := (*current).(type) {
			case *consCell:
				cons.SetCdr(newValue)
				current = &newValue
			case *nilCell:
				*current = newValue
			default:
				panic("Unexpected cell type in cons expression")
			}

		} else {
			break
		}
	}

	return cons
}
