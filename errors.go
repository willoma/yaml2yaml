package yaml2yaml

import (
	"errors"
	"fmt"
	"strings"

	"go.yaml.in/yaml/v4"
)

// TODO: Reduce usage of ErrInvalid, create more specific errors (for instance, I created ErrNotMapping for this.)

// Errors.
var (
	ErrAlreadyDefined = errors.New("already defined")
	ErrInvalid        = errors.New("invalid")
	ErrNotBoolean     = errors.New("not boolean")
	ErrNotComparable  = errors.New("not comparable")
	ErrNotDefined     = errors.New("not defined")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotMapping     = errors.New("not a mapping")
	ErrSyntaxError    = errors.New("syntax error")
	ErrUnsupported    = errors.New("unsupported")
)

// TODO Mapping: also have the index

type locatedError struct {
	path string
	node *yaml.Node
	err  error
}

func (e *locatedError) Error() string {
	nodeHasLoc := e.node != nil && e.node.Line > 0 && e.node.Column > 0

	switch {
	case nodeHasLoc && e.path != "":
		return fmt.Sprintf("[%d:%d] %s: %s", e.node.Line, e.node.Column, e.path, e.err)
	case nodeHasLoc:
		return fmt.Sprintf("[%d:%d]: %s", e.node.Line, e.node.Column, e.err)
	case e.path != "":
		return fmt.Sprintf("%s: %s", e.path, e.err)
	default:
		return e.err.Error()
	}
}

func (e *locatedError) Unwrap() error {
	return e.err
}

func (n *node) locatedError(err error) *locatedError {
	return &locatedError{path: n.path, node: n.Node, err: err}
}

func (m mappingChild) valueLocatedError(err error) *locatedError {
	if m.value == nil {
		return &locatedError{err: err}
	}

	return &locatedError{path: m.path, node: m.value, err: err}
}

type multiErrors []error

func (e multiErrors) Error() string {
	var buf strings.Builder

	buf.WriteString(e[0].Error())

	for _, err := range e[1:] {
		buf.WriteString(", ")
		buf.WriteString(err.Error())
	}

	return buf.String()
}

func joinErrors(errs ...error) error {
	result := make([]error, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			result = append(result, e)
		}
	}

	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	default:
		return multiErrors(result)
	}
}
