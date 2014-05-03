package main

import (
	"log"
)

type operation byte

const (
	opNil operation = iota
	opLdc
	opLdf
	opLd
	opAp
	opRet
	opSel
	opJoin
	opCons
	opCar
	opCdr
	opAdd
	opSub
	opMul
	opDiv
	opNeg
	opEq
	opStop
)

func run(code cell) cell {
	var op *opCell

	var stack cell = newNilCell()
	var env cell = newNilCell()
	var dump cell = newNilCell()

	running := true

	var jumpTable [opStop + 1]func()

	jumpTable[opNil] = func() {
		push(&stack, newNilCell())
	}

	jumpTable[opLdc] = func() {
		push(&stack, op.Data())
	}

	jumpTable[opLdf] = func() {
		push(&stack, newConsCell(op.Data(), env))
	}

	jumpTable[opLd] = func() {
		depth := op.Data().(*consCell).Car().(*fixNumCell).Value()
		position := op.Data().(*consCell).Cdr().(*fixNumCell).Value()
		push(&stack, env.(list).At(depth).(list).At(position))
	}

	jumpTable[opAp] = func() {
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

	jumpTable[opRet] = func() {
		returnValue := pop(&stack)

		// Restore old environment
		code = pop(&dump)
		env = pop(&dump)
		stack = pop(&dump)

		push(&stack, returnValue)
	}

	jumpTable[opSel] = func() {
		data := op.Data().(*consCell)

		push(&dump, code)

		switch pop(&stack).(type) {
		case *nilCell:
			code = data.Cdr()
		default:
			code = data.Car()
		}
	}

	jumpTable[opJoin] = func() {
		code = pop(&dump)
	}

	jumpTable[opCons] = func() {
		car := pop(&stack)
		cdr := pop(&stack)
		push(&stack, newConsCell(car, cdr))
	}

	jumpTable[opCar] = func() {
		push(&stack, pop(&stack).(*consCell).Car())
	}

	jumpTable[opCdr] = func() {
		push(&stack, pop(&stack).(*consCell).Cdr())
	}

	jumpTable[opAdd] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs+rhs))
	}

	jumpTable[opSub] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs-rhs))
	}

	jumpTable[opMul] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs*rhs))
	}

	jumpTable[opDiv] = func() {
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs/rhs))
	}

	jumpTable[opNeg] = func() {
		push(&stack, newFixNumCell(-pop(&stack).(*fixNumCell).Value()))
	}

	jumpTable[opEq] = func() {
		lhs := pop(&stack)
		rhs := pop(&stack)

		if lhs.Equal(rhs) {
			push(&stack, newTrueCell())
		} else {
			push(&stack, newNilCell())
		}
	}

	jumpTable[opStop] = func() {
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
