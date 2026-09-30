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

	if e.err == nil {
		return ""
	}

	//nolint:errorlint // We need to check precisely this err, not its wrapped errors recursively.
	if multiErr, ok := e.err.(*MultiError); ok {
		return e.errorMulti(path, multiErr)
	}

	return e.errorSingle(path, e.err)
}

func (e *located) errorMulti(path string, multi *MultiError) string {
	switch len(multi.errs) {
	case 0:
		return ""
	case 1:
		return e.errorSingle(path, multi.errs[0])
	default:
		var str strings.Builder

		str.WriteString(e.prefix(path))

		for _, err := range multi.errs {
			str.WriteString("\n  ")
			str.WriteString(strings.ReplaceAll(err.Error(), "\n", "\n  "))
		}

		return str.String()
	}
}

func (e *located) errorSingle(path string, err error) string {
	//nolint:errorlint // We need to check precisely this err, not its wrapped errors recursively.
	if err, ok := err.(*located); ok {
		return err.error(path)
	}

	return e.prefix(path) + " " + err.Error()
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
//
// If the error is nil or a [MultiError] with no errors, Locate returns nil.
func Locate(key string, line, column int, err error) error {
	if err == nil {
		return nil
	}

	//nolint:errorlint // We need to check precisely this err, not its wrapped errors recursively.
	if mErr, ok := err.(*MultiError); ok && mErr.IsNil() {
		return nil
	}

	return &located{key: key, line: line, column: column, err: err}
}
