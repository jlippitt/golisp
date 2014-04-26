package main

import (
	"io"
	"log"
)

type byteStream interface {
	io.Reader
	io.ByteReader
}

type vm struct {
	S cell
	E cell
	C byteStream
	D cell
}

type operation byte

const (
	OP_NIL operation = 0x00
	OP_LDC operation = 0x01
	OP_ADD operation = 0x02
)

func run(code byteStream) cell {
	var op byte
	var err error
	var opNil, opLdc, opAdd func()

	var stack cell = newNilCell()
	//var env cell = newNilCell()
	//var dump cell = newNilCell()

	opNil = func() {
		push(&stack, newNilCell())
		log.Print("NIL")
	}

	opLdc = func() {
		var bytes [8]byte
		var intVal int64 = 0

		if _, err := code.Read(bytes[:]); err != nil {
			panic(err)
		}

		for _, byte := range bytes {
			intVal = intVal << 8
			intVal += int64(byte)
		}

		push(&stack, newFixNumCell(intVal))

		log.Printf("LDC %d", intVal)
	}

	opAdd = func() {
		rhs := pop(&stack).(*fixNumCell).Value()
		lhs := pop(&stack).(*fixNumCell).Value()
		push(&stack, newFixNumCell(lhs+rhs))
		log.Printf("ADD")
	}

	jumpTable := []func(){
		opNil,
		opLdc,
		opAdd,
	}

	for {
		op, err = code.ReadByte()

		if err != nil {
			break
		}

		jumpTable[op]()
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
