package main

import "strconv"

type cellType byte

type cell interface{}

// NIL

type nilCell struct{}

var theNilCell nilCell = nilCell{}

func newNilCell() *nilCell {
	return &theNilCell
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

// PUSH AND POP

func push(list *cell, value cell) {
	*list = newConsCell(value, *list)
}

func pop(list *cell) cell {
	cons := (*list).(*consCell)
	value := cons.Car()
	*list = cons.Cdr()
	return value
}

// DUMP

func dump(cell cell) string {
	switch cell := cell.(type) {
	case *nilCell:
		return "()"
	case *consCell:
		return "(" + dump(cell.Car()) + " . " + dump(cell.Cdr()) + ")"
	case *symbolCell:
		return cell.Value()
	case *fixNumCell:
		return strconv.FormatInt(cell.Value(), 10)
	default:
		return "<unknown>"
	}
}
