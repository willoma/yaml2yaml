package errors

import (
	"fmt"
	"strings"
)

type located struct {
	key    string
	line   int
	column int
	err    error
}

func (e *located) Error() string { return e.error("") }
func (e *located) Unwrap() error { return e.err }

func (e *located) error(parentPath string) string {
	path := e.combinePath(parentPath)

	//nolint:errorlint // We need to check precisely this err, not its wrapped errors recursively.
	if err, ok := e.err.(*located); ok {
		return err.error(path)
	}

	msg := e.err.Error()

	if strings.Contains(msg, "\n") {
		return e.prefix(path) + "\n  " + strings.ReplaceAll(msg, "\n", "\n  ")
	}

	return e.prefix(path) + " " + msg
}

func (e *located) combinePath(parent string) string {
	if parent == "" {
		return e.key
	}

	if e.key == "" {
		return parent + "..."
	}

	if e.key[0] == '[' {
		return parent + e.key
	}

	return parent + "." + e.key
}

func (e *located) prefix(path string) string {
	hasPosition := e.line > 0 && e.column > 0

	switch {
	case hasPosition && path != "":
		return fmt.Sprintf("[%d:%d] %s:", e.line, e.column, path)
	case hasPosition:
		return fmt.Sprintf("[%d:%d]:", e.line, e.column)
	case path != "":
		return path + ":"
	default:
		return ""
	}
}

// Locate wraps an error with location information.
func Locate(key string, pos [2]int, err error) error {
	return &located{key: key, line: pos[0], column: pos[1], err: err}
}
