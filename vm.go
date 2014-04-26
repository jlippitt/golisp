package main

import (
	"log"
)

type operation byte

const (
	OP_NIL  operation = 0x00
	OP_LDC  operation = 0x01
	OP_ADD  operation = 0x02
	OP_HALT operation = 0x03
)

func run(code cell) cell {
	var op *opCell
	var opNil, opLdc, opAdd, opHalt func()

	var stack cell = newNilCell()
	//var env cell = newNilCell()
	//var dump cell = newNilCell()

	running := true

	opNil = func() {
		push(&stack, newNilCell())
		log.Print("NIL")
	}

	opLdc = func() {
		push(&stack, op.Data())
		log.Printf("LDC %s", dump(op.Data()))
	}

	opAdd = func() {
		rhs := pop(&stack).(*fixNumCell).Value()
		lhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs+rhs))
		log.Printf("ADD")
	}

	opHalt = func() {
		running = false
		log.Printf("HALT")
	}

	jumpTable := []func(){
		opNil,
		opLdc,
		opAdd,
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
