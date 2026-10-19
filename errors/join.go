package errors

import (
	"errors"
	"strings"
)

type JoinedError []error

func (e JoinedError) Append(err error) JoinedError {
	if errors.Is(err, ErrEmpty) {
		return e
	}

	return append(e, err)
}

func (e JoinedError) Error() string {
	var result strings.Builder

	result.WriteString(e[0].Error())

	for _, err := range e[1:] {
		result.WriteByte('\n')
		result.WriteString(err.Error())
	}

	return result.String()
}

func (e JoinedError) Compact() error {
	switch len(e) {
	case 0:
		return nil
	case 1:
		return e[0]
	default:
		return e
	}
}
