package main

import (
	"log"
)

type operation byte

const (
	OP_NIL operation = iota
	OP_LDC
	OP_CONS
	OP_AP
	OP_SEL
	OP_JOIN
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_EQ
	OP_HALT
)

func run(code cell) cell {
	var op *opCell

	var stack cell = newNilCell()
	//var env cell = newNilCell()
	var dump cell = newNilCell()

	running := true

	var jumpTable [OP_HALT + 1]func()

	jumpTable[OP_NIL] = func() {
		push(&stack, newNilCell())
	}

	jumpTable[OP_LDC] = func() {
		push(&stack, op.Data())
	}

	jumpTable[OP_CONS] = func() {
		car := pop(&stack)
		cdr := pop(&stack)
		push(&stack, newConsCell(car, cdr))
	}

	jumpTable[OP_AP] = func() {
		function := pop(&stack).(*functionCell)
		args := pop(&stack).(list)
		push(&stack, function.Call(args))
	}

	jumpTable[OP_SEL] = func() {
		data := op.Data().(*consCell)

		push(&dump, code)

		switch pop(&stack).(type) {
		case *nilCell:
			code = data.Cdr()
		default:
			code = data.Car()
		}
	}

	jumpTable[OP_JOIN] = func() {
		code = pop(&dump)
	}

	jumpTable[OP_ADD] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs+rhs))
	}

	jumpTable[OP_SUB] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs-rhs))
	}

	jumpTable[OP_MUL] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs*rhs))
	}

	jumpTable[OP_DIV] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs/rhs))
	}

	jumpTable[OP_EQ] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()

		if lhs == rhs {
			push(&stack, newTrueCell())
		} else {
			push(&stack, newNilCell())
		}
	}

	jumpTable[OP_HALT] = func() {
		running = false
	}

	for running {
		op = pop(&code).(*opCell)
		log.Print(op)
		jumpTable[op.Operation()]()
	}

	switch returnValue := stack.(type) {
	case *consCell:
		return returnValue.Car()
	case *nilCell:
		return nil
	default:
		panic("Unexpected cell type on stack")
	}
}
