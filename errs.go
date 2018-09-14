package errs

import (
	"fmt"

	"github.com/juju/errors"
)

type Error struct {
	Err error
}

func (this *Error) Error() string {
	return errors.ErrorStack(this.Err)
}

func New(msg string, args ...interface{}) error {
	return &Error{errors.New(fmt.Sprintf(msg, args...))}
}

func Append(err error, msg string, args ...interface{}) error {
	return &Error{errors.Annotatef(err, msg, args...)}
}
