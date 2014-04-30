package main

import (
	"fmt"
	"strconv"
)

type cellType byte

type cell interface {
	fmt.Stringer
}

type list interface {
	Current() cell
	Next() list
	IsNil() bool
}

// NIL

type nilCell struct{}

var theNilCell nilCell = nilCell{}

func newNilCell() *nilCell {
	return &theNilCell
}

func (self *nilCell) Current() cell {
	panic("Tried to get current value of nil")
}

func (self *nilCell) Next() list {
	panic("Tried to get next value of nil")
}

func (self *nilCell) IsNil() bool {
	return true
}

func (self *nilCell) String() string {
	return "()"
}

// CONS

type consCell struct {
	car cell
	cdr cell
}

func newConsCell(car cell, cdr cell) *consCell {
	return &consCell{car: car, cdr: cdr}
}

func (self *consCell) Car() cell {
	return self.car
}

func (self *consCell) SetCar(car cell) {
	self.car = car
}

func (self *consCell) Cdr() cell {
	return self.cdr
}

func (self *consCell) SetCdr(cdr cell) {
	self.cdr = cdr
}

func (self *consCell) Current() cell {
	return self.car
}

func (self *consCell) Next() list {
	return self.cdr.(list)
}

func (self *consCell) IsNil() bool {
	return false
}

func (self *consCell) String() string {
	return "(" + self.car.String() + " . " + self.cdr.String() + ")"
}

// OPCODE

type opCell struct {
	op   operation // See vm.go
	data cell
}

func newOpCell(op operation, data cell) *opCell {
	return &opCell{op: op, data: data}
}

func (self *opCell) Operation() operation {
	return self.op
}

func (self *opCell) Data() cell {
	return self.data
}

func (self *opCell) String() string {
	var output string

	switch self.op {
	case OP_NIL:
		output = "NIL"
	case OP_LDC:
		output = "LDC"
	case OP_CONS:
		output = "CONS"
	case OP_AP:
		output = "AP"
	case OP_HALT:
		output = "HALT"
	default:
		panic("Unknown opcode")
	}

	if self.data != nil {
		output += " " + self.data.String()
	}

	return output
}

// SYMBOL

type symbolCell struct {
	value string
}

func newSymbolCell(value string) *symbolCell {
	return &symbolCell{value: value}
}

func (self *symbolCell) Value() string {
	return self.value
}

func (self *symbolCell) String() string {
	return self.value
}

// FIXNUM

type fixNumCell struct {
	value int64
}

func newFixNumCell(value int64) *fixNumCell {
	return &fixNumCell{value: value}
}

func (self *fixNumCell) Value() int64 {
	return self.value
}

func (self *fixNumCell) String() string {
	return strconv.FormatInt(self.value, 10)
}

// BUILT-IN FUNCTION

type functionCell struct {
	name     string
	function func(list) cell
}

func newFunctionCell(name string, function func(list) cell) *functionCell {
	return &functionCell{name: name, function: function}
}

func (self *functionCell) Name() string {
	return self.name
}

func (self *functionCell) Call(args list) cell {
	return self.function(args)
}

func (self *functionCell) String() string {
	return "<BUILTIN:" + self.name + ">"
}

// PUSH AND POP

func push(list *cell, value cell) {
	*list = newConsCell(value, *list)
}

func pop(list *cell) cell {
	var value cell
	cons := (*list).(*consCell)
	value = cons.Car()
	*list = cons.Cdr()
	return value
}

func pushBack(it **cell, value cell) {
	var cdr cell = newConsCell(value, newNilCell())

	switch cons := (**it).(type) {
	case *consCell:
		cons.SetCdr(cdr)
		*it = &cdr
	case *nilCell:
		**it = cdr
	default:
		panic("Unexpected cell type in cons expression")
	}
}
