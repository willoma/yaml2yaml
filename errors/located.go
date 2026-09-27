package errors

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// Located returns an error that identifies the wrapped error with its source location.
func Located(path string, node *yaml.Node, err error) error {
	return &located{path: path, node: node, err: err}
}

type located struct {
	path string
	node *yaml.Node
	err  error
}

func (e *located) Error() string {
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

func (e *located) Unwrap() error {
	return e.err
}
