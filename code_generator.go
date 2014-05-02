package main

type codeWriter struct {
	code cell
	it   *cell
}

func newCodeWriter() *codeWriter {
	self := codeWriter{code: newNilCell()}
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

	function := node.(*consCell).Car()

	switch function := function.(type) {
	case *symbolCell:
		switch name := function.Value(); name {
		case "if":
			self.expandIf(node, args)
		case "+", "-", "*", "/", "=":
			self.expandOperator(node, name, args)
		default:
			// Named function call
		}
	case *consCell:
		self.expandFunctionCall(node, function, args)
	}
}

func (self *codeWriter) expandIf(node cell, args []cell) {
	if len(args) > 3 {
		panic("'if' expects no more than 3 arguments")
	} else if len(args) < 2 {
		panic("'if' expects at least 2 arguments")
	}

	lhs, rhs := newCodeWriter(), newCodeWriter()

	self.ExpandExpression(args[0])

	lhs.ExpandExpression(args[1])
	lhs.Write(OP_JOIN, nil)

	if len(args) > 2 {
		rhs.ExpandExpression(args[2])
	} else {
		rhs.Write(OP_NIL, nil)
	}

	rhs.Write(OP_JOIN, nil)

	self.Write(OP_SEL, newConsCell(lhs.Code(), rhs.Code()))
}

func (self *codeWriter) expandOperator(node cell, name string, args []cell) {
	if len(args) < 2 {
		panic("'" + name + "' expecets at least 2 arguments")
	}

	var op operation

	switch name {
	case "+":
		op = OP_ADD
	case "-":
		op = OP_SUB
	case "*":
		op = OP_MUL
	case "/":
		op = OP_DIV
	case "=":
		op = OP_EQ
	}

	// Push arguments to the stack in reverse order
	for i := len(args) - 1; i >= 0; i-- {
		self.ExpandExpression(args[i])
	}

	// Operator is applied from left to right
	for i := 0; i < len(args)-1; i++ {
		self.Write(op, nil)
	}
}

func (self *codeWriter) expandFunctionCall(node cell, function *consCell, args []cell) {
	// Push arguments to the stack in reverse order
	self.Write(OP_NIL, nil)

	for i := len(args) - 1; i >= 0; i-- {
		self.ExpandExpression(args[i])
		self.Write(OP_CONS, nil)
	}

	functionCode := newCodeWriter()
	functionCode.ExpandExpression(function)
	functionCode.Write(OP_RET, nil)

	self.Write(OP_LDF, functionCode.Code())
	self.Write(OP_AP, nil)
}

func generateCode(ast cell) cell {
	code := newCodeWriter()
	code.ExpandExpression(ast)
	code.Write(OP_HALT, nil)
	return code.Code()
}
