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

	byteCode := generateCode(parse(os.Args[1]))

	result := run(byteCode)

	fmt.Println("= " + dump(result))
}
