package main

type symbolTable map[string]cell

func newSymbolTable() *symbolTable {
	self := symbolTable{}

	self["+"] = newFunctionCell("+", func(args list) cell {
		total := args.Current().(*fixNumCell).Value()

		for args := args.Next(); !args.IsNil(); args = args.Next() {
			total += args.Current().(*fixNumCell).Value()
		}

		return newFixNumCell(total)
	})

	return &self
}

func (self *symbolTable) Get(symbol string) cell {
	return (*self)[symbol]
}
