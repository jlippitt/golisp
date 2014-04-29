package main

type symbolTable map[string]cell

func newSymbolTable() *symbolTable {
	self := symbolTable{}

	reduce := func(args list, callback func(total int64, value int64) int64) cell {
		total := args.Current().(*fixNumCell).Value()

		for args := args.Next(); !args.IsNil(); args = args.Next() {
			total = callback(total, args.Current().(*fixNumCell).Value())
		}

		return newFixNumCell(total)
	}

	self["+"] = newFunctionCell("+", func(args list) cell {
		return reduce(args, func(total int64, value int64) int64 {
			return total + value
		})
	})

	self["-"] = newFunctionCell("-", func(args list) cell {
		return reduce(args, func(total int64, value int64) int64 {
			return total - value
		})
	})

	self["*"] = newFunctionCell("*", func(args list) cell {
		return reduce(args, func(total int64, value int64) int64 {
			return total * value
		})
	})

	self["/"] = newFunctionCell("/", func(args list) cell {
		return reduce(args, func(total int64, value int64) int64 {
			return total / value
		})
	})

	return &self
}

func (self *symbolTable) Get(symbol string) cell {
	return (*self)[symbol]
}
