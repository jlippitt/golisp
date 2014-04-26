package main

import (
	"bytes"
	"encoding/binary"
)

func generateCode(ast cell) *bytes.Reader {
	var buf bytes.Buffer
	var expandExpression, expandFunctionCall func(cell)

	putOp := func(op operation) {
		if err := buf.WriteByte(byte(op)); err != nil {
			panic(err)
		}
	}

	putData := func(value interface{}) {
		if err := binary.Write(&buf, binary.BigEndian, value); err != nil {
			panic(err)
		}
	}

	expandExpression = func(node cell) {
		switch node := node.(type) {
		case *nilCell:
			putOp(OP_NIL)
		case *consCell:
			expandFunctionCall(node)
		case *fixNumCell:
			putOp(OP_LDC)
			putData(node.Value())
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

		// Arguments must be pushed on to the stack in reverse order
		for i := len(args) - 1; i >= 0; i-- {
			expandExpression(args[i])
		}

		putOp(OP_ADD)
	}

	expandExpression(ast)

	return bytes.NewReader(buf.Bytes())
}
