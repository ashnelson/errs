package errs

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

const noTraceMsg = "Unable to retrieve error stack trace"

var (
	callDepth int
	formatter Formatter
)

func init() {
	// callDepth defaults to 2 in order to ignore the stack trace from this package
	callDepth = 2
	formatter = &DefaultFormatter{}
}

// tracedError is a wrapper for the standard error object which keeps track of
// the previous error that was encountered
type tracedError struct {
	errDetails  string
	previousErr error
}

// Error returns the error details as a string
func (this *tracedError) Error() string {
	return this.errDetails
}

// New returns a new error with stack trace info added as part of the details
func New(msg string, args ...interface{}) error {
	return &tracedError{
		errDetails:  getFormattedStackTrace(msg, args...),
		previousErr: nil,
	}
}

// Append adds the "prevErr" details to a new error and also includes the stack
// trace info as part of the details of the new error. If "prevErr" is nil, a
// new error is returned with just the stack trace info added to the details.
func Append(prevErr error, msg string, args ...interface{}) error {
	if prevErr == nil {
		return New(msg, args...)
	}

	return &tracedError{
		errDetails:  fmt.Sprintf("%s;\n%s", prevErr, getFormattedStackTrace(msg, args...)),
		previousErr: prevErr,
	}
}

// SetCallDepth sets the call depth for retrieving the formatted stack trace
// which uses runtime.Caller(). This can be helpful if this package were to be
// wrapped by another package to allow the call depth to be changed to omit
// the wrapper package from the stack trace.
func SetCallDepth(depth int) {
	callDepth = depth
}

// SetFormatter sets the string formatter that formats each individual error
// line in the stack trace. The default formatter formats each line as short as
// possible and looks similar to the following:
//	filename1.go:func1():123: message1
//	filename2.go:func2():456: message2
//	filename3.go:func3():789: message3
func SetFormatter(fmtr Formatter) {
	formatter = fmtr
}

// getFormattedStackTrace gets the caller information (file, function, and line
// number) and returns it as formatted string
func getFormattedStackTrace(msg string, args ...interface{}) string {
	// Get caller informtion
	prgrmCntr, tracedFile, tracedLnNbr, traceReceived := runtime.Caller(callDepth)
	if !traceReceived || tracedFile == "" || tracedLnNbr < 1 {
		return noTraceMsg
	}

	// Should never be nil but check to be safe
	funcForPC := runtime.FuncForPC(prgrmCntr)
	if funcForPC == nil {
		return noTraceMsg
	}

	tracedFile = path.Base(tracedFile)
	tracedFunc := path.Base(funcForPC.Name())

	// tracedFunc generally looks like "fileName.functionName"; attempt to get
	// just the function name
	if splitFuncName := strings.Split(tracedFunc, "."); len(splitFuncName) == 2 {
		tracedFunc = splitFuncName[1]
	}

	return formatter.Format(tracedFile, tracedFunc, tracedLnNbr, msg, args...)
}
