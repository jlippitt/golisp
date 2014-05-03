package main

import "fmt"

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
		self.Write(opNil, nil)
	case *trueCell:
		self.Write(opLdc, newTrueCell())
	case *consCell:
		self.expandForm(node)
	case *fixNumCell:
		self.Write(opLdc, newFixNumCell(node.Value()))
	case *symbolCell:
		location := self.st.Locate(node.Value())
		self.Write(opLd, newConsCell(
			newFixNumCell(location.Depth),
			newFixNumCell(location.Position),
		))
	default:
		panic("Unexpected node type in expression")
	}
}

func (self *codeWriter) expandForm(node *consCell) {
	function := node.Car()
	args := node.Cdr().(list).Slice()

	switch function := function.(type) {
	case *symbolCell:
		switch name := function.Value(); name {
		case "if":
			self.expandIf(args)
			return
		case "fn":
			self.expandAnonymousFunction(args)
			return
		case "cons":
			self.expandCons(args)
			return
		case "car", "cdr", "atom", "+", "-", "*", "/", "=":
			self.expandOperator(name, args)
			return
		}
	}

	self.expandFunctionCall(function, args)
}

func (self *codeWriter) expandIf(args []cell) {
	if len(args) > 3 {
		panic("'if' expects no more than 3 arguments")
	} else if len(args) < 2 {
		panic("'if' expects at least 2 arguments")
	}

	lhs, rhs := newCodeWriter(self.st), newCodeWriter(self.st)

	self.ExpandExpression(args[0])

	lhs.ExpandExpression(args[1])
	lhs.Write(opJoin, nil)

	if len(args) > 2 {
		rhs.ExpandExpression(args[2])
	} else {
		rhs.Write(opNil, nil)
	}

	rhs.Write(opJoin, nil)

	self.Write(opSel, newConsCell(lhs.Code(), rhs.Code()))
}

func (self *codeWriter) expandAnonymousFunction(args []cell) {
	body := newCodeWriter(self.st)

	self.st.UpLevel()

	for _, param := range args[0].(list).Slice() {
		self.st.Register(param.(*symbolCell).Value())
	}

	body.ExpandExpression(args[1])

	self.st.DownLevel()

	body.Write(opRet, nil)
	self.Write(opLdf, body.Code())
}

func (self *codeWriter) expandCons(args []cell) {
	self.Write(opNil, nil)

	// Push arguments to the stack in reverse order
	for i := len(args) - 1; i >= 0; i-- {
		self.ExpandExpression(args[i])
		self.Write(opCons, nil)
	}
}

func (self *codeWriter) expandOperator(name string, args []cell) {
	var op operation
	var minArgs, maxArgs int

	switch name {
	case "car":
		op = opCar
		minArgs, maxArgs = 1, 1
	case "cdr":
		op = opCdr
		minArgs, maxArgs = 1, 1
	case "atom":
		op = opAtom
		minArgs, maxArgs = 1, 1
	case "+":
		op = opAdd
		minArgs = 2
	case "-":
		if len(args) > 1 {
			op = opSub
		} else {
			op = opNeg
		}
		minArgs = 1
	case "*":
		op = opMul
		minArgs = 2
	case "/":
		op = opDiv
		minArgs = 2
	case "=":
		op = opEq
		minArgs = 2
	}

	if minArgs != 0 && len(args) < minArgs {
		panic(fmt.Sprintf("'%s' expects at least %d arguments", name, minArgs))
	} else if maxArgs != 0 && len(args) > maxArgs {
		panic(fmt.Sprintf("'%s' expects no more than %d arguments", name, maxArgs))
	}

	// Push arguments to the stack in reverse order
	for i := len(args) - 1; i >= 0; i-- {
		self.ExpandExpression(args[i])
	}

	// Operator is applied from left to right
	if len(args) > 1 {
		for i := 0; i < len(args)-1; i++ {
			self.Write(op, nil)
		}
	} else {
		self.Write(op, nil)
	}
}

func (self *codeWriter) expandFunctionCall(function cell, args []cell) {
	self.expandCons(args)
	self.ExpandExpression(function)
	self.Write(opAp, nil)
}

func generateCode(ast cell) cell {
	code := newCodeWriter(newSymbolTable())
	code.ExpandExpression(ast)
	code.Write(opStop, nil)
	return code.Code()
}
