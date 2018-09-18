# errs

The `errs` pacakge offers a seamless way to add a stack trace to any error 
object. The use of `errs` is very similar to the way the Go standard library
handles errors and the `errs` error object can be used just like the standard
libarary error object for use in logging or even printing to stdout or stderr.

## How to use
A simple example use case would look something like the following:
```Go
package main

import (
	"fmt"

	"errs"
)

func main() {
	err := testFunc1()
	fmt.Println(err)
}

func testFunc1() error {
	if err := testFunc2(); err != nil {
		return errs.Append(err, "Oh no! Another error occurred!")
	}
	return nil
}

func testFunc2() error {
	return errs.New("Something went wrong!")
}
```

The above sample code would output something like the following to stdout:
```
main.go:testFunc2():22: Something went wrong!;
main.go:testFunc1():16: Oh no! Another error occurred!
```

## Creating a custom Formatter
The stack trace formatting can also easily be changed by implementing new
`errs.Formatter` similar to the following:
```Go
type VerboseFormatter struct{}

func (this *VerboseFormatter) Format(tracedFile, tracedFunc string, tracedLnNbr int, msg string) string {
	return fmt.Sprintf("The error occurred in %q in function %q on line %d;\n\tDetails: %s", tracedFile, tracedFunc, tracedLnNbr, msg)
}

func main() {
	// Make the errs package use the new Formatter
	errs.SetFormatter(&VerboseFormatter{})
	...
}
```

The output of the above sample code, with the new `VerboseFormatter`, would then look like the following:
```
The error occurred in "main.go" in function "testFunc2" on line 30;
	Details: Something went wrong!;
The error occurred in "main.go" in function "testFunc1" on line 24;
	Details: Oh no! Another error occurred!
```