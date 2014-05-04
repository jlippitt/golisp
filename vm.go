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
	opSt
	opAp
	opRet
	opSel
	opJoin
	opCons
	opCar
	opCdr
	opAtom
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
	var env cell = newConsCell(newNilCell(), newNilCell())
	var dump cell = newNilCell()

	running := true

	var jumpTable [opStop + 1]func()

	jumpTable[opNil] = func() {
		// Push a nil cell on to the stack
		push(&stack, newNilCell())
	}

	jumpTable[opLdc] = func() {
		// Push a constant on to the stack
		push(&stack, op.Data())
	}

	jumpTable[opLdf] = func() {
		// Push a function (closure) on to the stack
		push(&stack, newConsCell(op.Data(), env))
	}

	jumpTable[opLd] = func() {
		// Load the value of a variable
		depth := op.Data().(*consCell).Car().(*fixNumCell).Value()
		position := op.Data().(*consCell).Cdr().(*fixNumCell).Value()
		push(&stack, env.(list).At(depth).(list).At(position))
	}

	jumpTable[opSt] = func() {
		// Store a value as a variable
		env := env.(*consCell)

		var scope cell = env.Car()

		value := newConsCell(pop(&stack), newNilCell())

		switch scope.(type) {
		case *consCell:
			var last list

			for it := scope.(list); !it.IsNil(); it = it.Next() {
				last = it
			}

			last.(*consCell).SetCdr(value)

		case *nilCell:
			env.SetCar(value)
		}
	}

	jumpTable[opAp] = func() {
		// Apply a function
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
		// Return from a function
		returnValue := pop(&stack)

		// Restore old environment
		code = pop(&dump)
		env = pop(&dump)
		stack = pop(&dump)

		push(&stack, returnValue)
	}

	jumpTable[opSel] = func() {
		// Choose between two code paths
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
		// Join the main code path again
		code = pop(&dump)
	}

	jumpTable[opCons] = func() {
		// Create a cons cell
		car := pop(&stack)
		cdr := pop(&stack)
		push(&stack, newConsCell(car, cdr))
	}

	jumpTable[opCar] = func() {
		// Get the head (car) of a cons cell
		push(&stack, pop(&stack).(*consCell).Car())
	}

	jumpTable[opCdr] = func() {
		// Get the tail (cdr) of a cons cell
		push(&stack, pop(&stack).(*consCell).Cdr())
	}

	jumpTable[opAtom] = func() {
		// Is the cell an atom (i.e. non-list cell)?
		switch pop(&stack).(type) {
		case list:
			push(&stack, newNilCell())
		default:
			push(&stack, newTrueCell())
		}
	}

	jumpTable[opAdd] = func() {
		// Add two integers
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs+rhs))
	}

	jumpTable[opSub] = func() {
		// Subtract two integers
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs-rhs))
	}

	jumpTable[opMul] = func() {
		// Multiply two integers
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs*rhs))
	}

	jumpTable[opDiv] = func() {
		// Divide two integers
		lhs := pop(&stack).(*fixNumCell).Value()
		rhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs/rhs))
	}

	jumpTable[opNeg] = func() {
		// Negate an integer
		push(&stack, newFixNumCell(-pop(&stack).(*fixNumCell).Value()))
	}

	jumpTable[opEq] = func() {
		// Are two values equal?
		lhs := pop(&stack)
		rhs := pop(&stack)

		if lhs.Equal(rhs) {
			push(&stack, newTrueCell())
		} else {
			push(&stack, newNilCell())
		}
	}

	jumpTable[opStop] = func() {
		// Stop execution
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
