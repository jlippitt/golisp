package main

type cellType byte

type cell interface{}

// NIL

type nilCell struct{}

func newNilCell() *nilCell {
	return &nilCell{}
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

func (self *consCell) Cdr() cell {
	return self.cdr
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

func push(list cell, value cell) cell {
	return newConsCell(value, list)
}

func pop(list cell) (cell, cell) {
	cons := list.(*consCell)
	return cons.Cdr(), cons.Car()
}
