package main

type codeWriter struct {
	code cell
	it   *cell
	st   *symbolTable
}

func newCodeWriter(st *symbolTable) *codeWriter {
	self := codeWriter{code: newNilCell(), st: st}
	self.it = &self.code
	return &self
}

func (self *codeWriter) Write(op operation, data cell) {
	pushBack(&self.it, newOpCell(op, data))
}

func (self *codeWriter) Code() cell {
	return self.code
}

func (self *codeWriter) ExpandExpression(node cell) {
	switch node := node.(type) {
	case *nilCell:
		self.Write(OP_NIL, nil)
	case *consCell:
		self.expandForm(node)
	case *fixNumCell:
		self.Write(OP_LDC, newFixNumCell(node.Value()))
	default:
		panic("Unexpected node type in expression")
	}
}

func (self *codeWriter) expandForm(node cell) {
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

	functionName := node.(*consCell).Car().(*symbolCell).Value()

	if functionName == "if" {
		if len(args) > 3 {
			panic("'if' expects no more than 3 arguments")
		} else if len(args) < 2 {
			panic("'if' expects at least 2 arguments")
		}

		lhs, rhs := newCodeWriter(self.st), newCodeWriter(self.st)

		self.ExpandExpression(args[0])

		lhs.ExpandExpression(args[1])
		lhs.Write(OP_JOIN, nil)

		if len(args) > 2 {
			rhs.ExpandExpression(args[2])
		}

		rhs.Write(OP_JOIN, nil)

		self.Write(OP_SEL, newConsCell(lhs.Code(), rhs.Code()))

	} else {
		// Normal function call
		self.Write(OP_NIL, nil)

		// Arguments must be pushed on to the stack in reverse order
		for i := len(args) - 1; i >= 0; i-- {
			self.ExpandExpression(args[i])
			self.Write(OP_CONS, nil)
		}

		self.Write(OP_LDC, self.st.Get(functionName))
		self.Write(OP_AP, nil)
	}
}

func generateCode(ast cell, st *symbolTable) cell {
	code := newCodeWriter(st)
	code.ExpandExpression(ast)
	code.Write(OP_HALT, nil)
	return code.Code()
}
