package errs

import (
	"encoding/json"
	"fmt"
)

// Formatter specifies exactly how each paramter should be formatted as a string.
// The format can be as terse or verbose as desired or even omit certain
// parameters if they aren't needed.
// For examples on writing a formatter, see DefaultFormatter and JSONFormatter.
type Formatter interface {
	Format(tracedFile, tracedFunc string, tracedLnNbr int, msg string, args ...interface{}) string
}

// DefaultFormatter formats the error as a simple string like:
//	filename1.go:func1():123: message1
//	filename2.go:func2():456: message2
//	filename3.go:func3():789: message3
type DefaultFormatter struct{}

var _ Formatter = &DefaultFormatter{}

func (this *DefaultFormatter) Format(tracedFile, tracedFunc string, tracedLnNbr int, msg string, args ...interface{}) string {
	msg = fmt.Sprintf(msg, args...)
	return fmt.Sprintf("%s:%s():%d: %s", tracedFile, tracedFunc, tracedLnNbr, msg)
}

// JSONFormatter formats the error as a simple string like:
//	{"File":"filename1.go","Func":"func1()","LineNbr":123,"Details":"message1"}
//	{"File":"filename2.go","Func":"func2()","LineNbr":456,"Details":"message2"}
//	{"File":"filename3.go","Func":"func3()","LineNbr":789,"Details":"message3"}
type JSONFormatter struct{}

var _ Formatter = &JSONFormatter{}

func (this *JSONFormatter) Format(tracedFile, tracedFunc string, tracedLnNbr int, msg string, args ...interface{}) string {
	type Fmtr struct {
		File    string
		Func    string
		LineNbr int
		Details string
	}
	jsonFmtr := &Fmtr{
		File:    tracedFile,
		Func:    tracedFunc + "()",
		LineNbr: tracedLnNbr,
		Details: fmt.Sprintf(msg, args...),
	}

	msgJSON, err := json.Marshal(jsonFmtr)
	if err != nil {
		return noTraceMsg + "Details: " + err.Error()
	}

	return string(msgJSON)
}
