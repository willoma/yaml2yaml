package yaml2yaml

import (
	"errors"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// Errors.
var (
	ErrAlreadyDefined  = errors.New("already defined")
	ErrCannotMerge     = errors.New("cannot merge")
	ErrInvalid         = errors.New("invalid")
	ErrNotBoolean      = errors.New("not boolean")
	ErrNotComparable   = errors.New("not comparable")
	ErrNotDefined      = errors.New("not defined")
	ErrNotImplemented  = errors.New("not implemented")
	ErrNotMapping      = errors.New("not a mapping")
	ErrSyntaxError     = errors.New("syntax error")
	ErrUnexpectedValue = errors.New("unexpected value")
	ErrUnsupported     = errors.New("unsupported")
)

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

func (m mappingChild) locatedError(err error) error {
	return &locatedError{path: m.path, node: m.value, err: err}
}
