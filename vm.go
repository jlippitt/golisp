package main

import (
	"log"
)

type operation byte

const (
	OP_NIL operation = iota
	OP_LDC
	OP_LDF
	OP_LD
	OP_CONS
	OP_AP
	OP_RET
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
	var env cell = newNilCell()
	var dump cell = newNilCell()

	running := true

	var jumpTable [OP_HALT + 1]func()

	jumpTable[OP_NIL] = func() {
		push(&stack, newNilCell())
	}

	jumpTable[OP_LDC] = func() {
		push(&stack, op.Data())
	}

	jumpTable[OP_LDF] = func() {
		push(&stack, newConsCell(op.Data(), env))
	}

	jumpTable[OP_LD] = func() {
		depth := op.Data().(*consCell).Car().(*fixNumCell).Value()
		position := op.Data().(*consCell).Cdr().(*fixNumCell).Value()
		push(&stack, env.(list).At(depth).(list).At(position))
	}

	jumpTable[OP_CONS] = func() {
		car := pop(&stack)
		cdr := pop(&stack)
		push(&stack, newConsCell(car, cdr))
	}

	jumpTable[OP_AP] = func() {
		function := pop(&stack).(*consCell)
		args := pop(&stack)

		// Save the old environment to the dump
		push(&dump, stack)
		push(&dump, env)
		push(&dump, code)

		// New environment
		stack = nil
		code = function.Car()
		env = function.Cdr()

		push(&env, args)
	}

	jumpTable[OP_RET] = func() {
		returnValue := pop(&stack)

		// Restore old environment
		code = pop(&dump)
		env = pop(&dump)
		stack = pop(&dump)

		push(&stack, returnValue)
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
		lhs := pop(&stack)
		rhs := pop(&stack)

		if lhs.Equal(rhs) {
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
