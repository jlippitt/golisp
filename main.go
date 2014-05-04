package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var result cell
	var code []byte
	var err error

	if len(os.Args) < 2 {
		fmt.Println("No source file specified")
		os.Exit(1)
	}

	if code, err = ioutil.ReadFile(os.Args[1]); err != nil {
		fmt.Println("Could not read source file")
		os.Exit(1)
	}

	result, err = Execute(string(code))

	if err != nil {
		fmt.Print("Error: " + err.Error())
		os.Exit(2)
	}

	fmt.Printf("= %s\n", result)
}
