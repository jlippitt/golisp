package main

func generateCode(ast cell, st *symbolTable) cell {
	var code cell = newNilCell()
	var it *cell = &code
	var expandExpression, expandFunctionCall func(cell)

	write := func(op operation, data cell) {
		pushBack(&it, newOpCell(op, data))
	}

	expandExpression = func(node cell) {
		switch node := node.(type) {
		case *nilCell:
			write(OP_NIL, nil)
		case *consCell:
			expandFunctionCall(node)
		case *fixNumCell:
			write(OP_LDC, newFixNumCell(node.Value()))
		default:
			panic("Unexpected node type in expression")
		}
	}

	expandFunctionCall = func(node cell) {
		// TODO: We'll always use add for now - add other stuff later
		var args []cell

		var current *consCell = node.(*consCell)

		// Traverse the list and get all arguments to the function
	Loop:
		for {
			switch cdr := current.Cdr().(type) {
			case *consCell:
				args = append(args, cdr.Car())
				current = cdr
			case *nilCell:
				break Loop
			default:
				panic("Unexpected node type in function call")
			}
		}

		write(OP_NIL, nil)

		// Arguments must be pushed on to the stack in reverse order
		for i := len(args) - 1; i >= 0; i-- {
			expandExpression(args[i])
			write(OP_CONS, nil)
		}

		functionName := node.(*consCell).Car().(*symbolCell).Value()

		write(OP_LDC, st.Get(functionName))
		write(OP_AP, nil)
	}

	expandExpression(ast)

	write(OP_HALT, nil)

	return code
}
