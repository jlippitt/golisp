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

	theVm := vm{
		S: newNilCell(),
		E: newNilCell(),
		C: code,
		D: newNilCell(),
	}

	jumpTable := []func(*vm){
		opNil,
		opLdc,
		opAdd,
	}

	for {
		op, err = theVm.C.ReadByte()

		if err != nil {
			break
		}

		jumpTable[op](&theVm)
	}

	switch returnValue := theVm.S.(type) {
	case *consCell:
		return returnValue.Car()
	case *nilCell:
		return nil
	default:
		panic("Unexpected cell type on stack")
	}
}

func opNil(vm *vm) {
	push(&vm.S, newNilCell())

	log.Print("NIL")
}

func opLdc(vm *vm) {
	var bytes [8]byte
	var intVal int64 = 0

	if _, err := vm.C.Read(bytes[:]); err != nil {
		panic(err)
	}

	for _, byte := range bytes {
		intVal = intVal << 8
		intVal += int64(byte)
	}

	push(&vm.S, newFixNumCell(intVal))

	log.Printf("LDC %d", intVal)
}

func opAdd(vm *vm) {
	rhs := pop(&vm.S).(*fixNumCell).Value()
	lhs := pop(&vm.S).(*fixNumCell).Value()

	push(&vm.S, newFixNumCell(lhs+rhs))

	log.Printf("ADD")
}
