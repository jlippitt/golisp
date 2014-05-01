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

	self["="] = newFunctionCell("=", func(args list) cell {
		lhs := args.Current()
		rhs := args.Next().Current()

		if lhs.Equal(rhs) {
			return newTrueCell()
		} else {
			return newNilCell()
		}
	})

	return &self
}

func (self *symbolTable) Get(symbol string) cell {
	return (*self)[symbol]
}
