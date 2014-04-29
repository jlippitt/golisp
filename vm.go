package main

import (
	"log"
)

type operation byte

const (
	OP_NIL  operation = 0x00
	OP_LDC  operation = 0x01
	OP_CONS operation = 0x02
	OP_AP   operation = 0x03
	OP_HALT operation = 0x04
)

func run(code cell) cell {
	var op *opCell
	var opNil, opLdc, opCons, opAp, opHalt func()

	var stack cell = newNilCell()
	//var env cell = newNilCell()
	//var dump cell = newNilCell()

	running := true

	opNil = func() {
		log.Print("NIL")
		push(&stack, newNilCell())
	}

	opLdc = func() {
		log.Printf("LDC %s", dump(op.Data()))
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

	opHalt = func() {
		log.Printf("HALT")
		running = false
	}

	jumpTable := []func(){
		opNil,
		opLdc,
		opCons,
		opAp,
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
