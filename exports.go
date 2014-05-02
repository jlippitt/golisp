package main

// When this is made a library, the exported API will be defined here
type Error string

func (self Error) Error() string {
	return string(self)
}

// TODO: 'cell' to be made public (i.e. 'Cell')?
func Execute(code string) (result cell, err error) {
	// TODO: Uncomment when we have better error handling
	/*defer func() {
		recoveredError := recover()

		if recoveredError != nil {
			switch recoveredError := recoveredError.(type) {
			case error:
				err = recoveredError
			case string:
				err = Error(recoveredError)
			default:
				err = Error("Unknown error")
			}
		}
	}()*/

	result = run(generateCode(parse(code)))
	return
}
