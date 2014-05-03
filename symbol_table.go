package main

import golist "container/list"

type symbolTable struct {
	list golist.List
}

type symbolMap struct {
	Positions map[string]int64
	Next      int64
}

type location struct {
	Depth    int64
	Position int64
}

func newSymbolTable() *symbolTable {
	var self symbolTable
	self.UpLevel()
	return &self
}

func (self *symbolTable) UpLevel() {
	self.list.PushFront(&symbolMap{Positions: make(map[string]int64)})
}

func (self *symbolTable) DownLevel() {
	self.list.Remove(self.list.Front())
}

func (self *symbolTable) Register(symbol string) {
	// TODO: Duplicate checking
	currentLevel := self.list.Front().Value.(*symbolMap)
	currentLevel.Positions[symbol] = currentLevel.Next
	currentLevel.Next++
}

func (self *symbolTable) Locate(symbol string) *location {
	var depth int64 = 0

	for el := self.list.Front(); el != nil; el = el.Next() {
		position, ok := el.Value.(*symbolMap).Positions[symbol]

		if ok {
			return &location{depth, position}
		}

		depth++
	}

	return nil
}
