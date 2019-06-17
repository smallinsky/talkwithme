package command

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCommandFormat = errors.New("invalid command format")
	ErrCmdNotSupported      = errors.New("command is not supported")
	ErrInvalidInput         = errors.New("invalid input")
	ErrEmptyLine            = errors.New("got empty line")
)

// Scanner validates and transforms user input to server a client command.
func Scan(input string) (interface{}, error) {
	return scanCommand(input)
}

func scanCommand(line string) (interface{}, error) {
	var (
		name string
		args []string
	)

	if len(line) == 0 {
		return nil, ErrEmptyLine
	}

	if line[0] == '/' {
		s := strings.Split(strings.TrimPrefix(line, "/"), " ")
		name = s[0]
		args = append(args, s[1:]...)
	} else {
		name = "textmessage"
		args = append(args, line)
	}
	return createCommand(name, args)
}

func createCommand(cmd string, args []string) (interface{}, error) {
	handler, ok := supportedCommands[cmd]
	if !ok {
		return nil, ErrCmdNotSupported
	}
	for _, validator := range handler.validators {
		if err := validator(args); err != nil {
			return nil, err
		}
	}
	return handler.bind(args), nil
}
