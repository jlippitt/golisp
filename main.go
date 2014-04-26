package main

import "log"

func main() {
	code := "(+ 6 5)"

	byteCode := generateCode(parse(code))

	result := run(byteCode)

	log.Print("= " + dump(result))
}
