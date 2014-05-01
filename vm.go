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
	OP_HALT
)

func run(code cell) cell {
	var op *opCell
	var opNil, opLdc, opCons, opAp, opSel, opJoin, opHalt func()

	var stack cell = newNilCell()
	//var env cell = newNilCell()
	var dump cell = newNilCell()

	running := true

	opNil = func() {
		log.Print("NIL")
		push(&stack, newNilCell())
	}

	opLdc = func() {
		log.Printf("LDC %s", op.Data())
		push(&stack, op.Data())
	}

	opCons = func() {
		log.Printf("CONS")
		car := pop(&stack)
		cdr := pop(&stack)
		push(&stack, newConsCell(car, cdr))
	}

	opAp = func() {
		log.Printf("AP")
		function := pop(&stack).(*functionCell)
		args := pop(&stack).(list)
		push(&stack, function.Call(args))
	}

	opSel = func() {
		log.Printf("SEL %s", op.Data())

		data := op.Data().(*consCell)

		push(&dump, code)

		switch pop(&stack).(type) {
		case *nilCell:
			code = data.Cdr()
		default:
			code = data.Car()
		}
	}

	opJoin = func() {
		log.Printf("JOIN")
		code = pop(&dump)
	}

	opHalt = func() {
		log.Printf("HALT")
		running = false
	}

	jumpTable := []func(){
		opNil,
		opLdc,
		opCons,
		opAp,
		opSel,
		opJoin,
		opHalt,
	}

	for running {
		op = pop(&code).(*opCell)
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
