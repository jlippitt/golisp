package main

import (
	"fmt"
	"strconv"
)

type cellType byte

type cell interface {
	fmt.Stringer
	Equal(cell) bool
}

type list interface {
	Current() cell
	Next() list
	IsNil() bool
	At(int64) cell
	Slice() []cell
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

func (self *nilCell) At(index int64) cell {
	panic("Tried to find value in empty list")
}

func (self *nilCell) Slice() []cell {
	return nil
}

func (self *nilCell) Equal(other cell) bool {
	switch other.(type) {
	case *nilCell:
		return true
	default:
		return false
	}
}

func (self *nilCell) String() string {
	return "()"
}

// TRUE

type trueCell struct{}

var theTrueCell trueCell = trueCell{}

func newTrueCell() *trueCell {
	return &theTrueCell
}

func (self *trueCell) Equal(other cell) bool {
	switch other.(type) {
	case *trueCell:
		return true
	default:
		return false
	}
}

func (self *trueCell) String() string {
	return "T"
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

func (self *consCell) At(index int64) cell {
	var it list = self

	for i := int64(0); i < index; i++ {
		it = it.Next()
	}

	return it.Current()
}

func (self *consCell) Slice() []cell {
	var elements []cell
	var it list

	for it = self; !it.IsNil(); it = it.Next() {
		elements = append(elements, it.Current())
	}

	return elements
}

func (self *consCell) IsNil() bool {
	return false
}

func (self *consCell) Equal(other cell) bool {
	switch other := other.(type) {
	case *consCell:
		return self.car.Equal(other.Car()) && self.cdr.Equal(other.Cdr())
	default:
		return false
	}
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

func (self *opCell) Equal(other cell) bool {
	switch other := other.(type) {
	case *opCell:
		if self.op == other.Operation() {
			return (self.data == nil && other.Data() == nil) ||
				self.data.Equal(other.Data())
		} else {
			return false
		}
	default:
		return false
	}
}

func (self *opCell) String() string {
	var output string

	switch self.op {
	case OP_NIL:
		output = "NIL"
	case OP_LDC:
		output = "LDC"
	case OP_LDF:
		output = "LDF"
	case OP_LD:
		output = "LD"
	case OP_CONS:
		output = "CONS"
	case OP_AP:
		output = "AP"
	case OP_RET:
		output = "RET"
	case OP_SEL:
		output = "SEL"
	case OP_JOIN:
		output = "JOIN"
	case OP_ADD:
		output = "ADD"
	case OP_SUB:
		output = "SUB"
	case OP_MUL:
		output = "MUL"
	case OP_DIV:
		output = "DIV"
	case OP_EQ:
		output = "EQ"
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

func (self *symbolCell) Equal(other cell) bool {
	switch other := other.(type) {
	case *symbolCell:
		return self.value == other.Value()
	default:
		return false
	}
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

func (self *fixNumCell) Equal(other cell) bool {
	switch other := other.(type) {
	case *fixNumCell:
		return self.value == other.Value()
	default:
		return false
	}
}

func (self *fixNumCell) String() string {
	return strconv.FormatInt(self.value, 10)
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
