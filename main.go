package main

import "log"

func main() {
	code := "(+ 6 5)"

	log.Print(dump(parse(code)))
}
