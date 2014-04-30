package main

import (
	"fmt"
	"os"
)

func main() {
	var result cell
	var err error

	if len(os.Args) < 2 {
		fmt.Println("No code to execute")
		os.Exit(1)
	}

	result, err = Execute(os.Args[1])

	if err != nil {
		fmt.Print("Error: " + err.Error())
		os.Exit(2)
	}

	fmt.Printf("= %s\n", result)
}
