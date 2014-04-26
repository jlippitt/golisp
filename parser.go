package main

func parse(input string) cell {
	var parseExpression, parseCons func() cell

	tokenizer := newTokenizer(input)

	parseExpression = func() cell {
		var value cell

		switch tokenizer.Type() {
		case TOK_OPEN:
			value = parseCons()
		case TOK_SYMBOL:
			value = newSymbolCell(tokenizer.StringValue())
		case TOK_FIXNUM:
			value = newFixNumCell(tokenizer.IntValue())
		default:
			panic("Unexpected token")
		}

		return value
	}

	parseCons = func() cell {
		var cons cell = newNilCell()
		var current *cell = &cons

		for {
			tokenizer.Next()

			if tokenizer.Type() != TOK_CLOSE {
				var newValue cell = newConsCell(parseExpression(), newNilCell())

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

	return parseExpression()
}
