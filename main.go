package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No code to execute")
		os.Exit(1)
	}

	st := newSymbolTable()

	byteCode := generateCode(parse(os.Args[1]), st)

	result := run(byteCode)

	fmt.Printf("= %s\n", result)
}
