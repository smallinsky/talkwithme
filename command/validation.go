package command

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidArgsCount = errors.New("invalid arguments count")
	ErrMaxLengthExeeded = errors.New("arg max length exceeded")
)

// numberOfArguments ensures that N arguments were provided.
func numberOfArguments(n int) validateFn {
	return func(args []string) error {
		if got := len(args); got != n {
			return errors.Wrapf(ErrInvalidArgsCount, "expect: %v got: %v", n, got)
		}
		return nil
	}
}

// argumentNMaxLength checks if length N argument  exceeds maxLength.
func argumentNMaxLength(argN, maxLength int) validateFn {
	return func(args []string) error {
		if got := len(args[argN]); got > maxLength {
			return errors.Wrapf(ErrMaxLengthExeeded, "arg[%v] max length: %v, got: %v", argN, maxLength, got)
		}
		return nil
	}
}
