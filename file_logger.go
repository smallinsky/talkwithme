package main

import (
	"fmt"
	"io"
	"os"
)

type fileLogger struct {
	file io.WriteCloser
}

func newFileLogger(path string) (*fileLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v", err)
	}
	return &fileLogger{
		file: f,
	}, nil
}

func (s *fileLogger) Close() error {
	return s.file.Close()
}

func (s *fileLogger) LogMessage(msg string) error {
	if len(msg) == 0 {
		return nil
	} else if msg[len(msg)-1] != '\n' {
		msg = fmt.Sprintf("%s\n", msg)
	}

	if _, err := s.file.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
