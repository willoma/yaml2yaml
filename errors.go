package yaml2yaml

import (
	oerrors "errors"

	"github.com/willoma/yaml2yaml/errors"
)

// Errors.
var (
	ErrAlreadyDefined  = oerrors.New("already defined")
	ErrCannotMerge     = oerrors.New("cannot merge")
	ErrInvalid         = oerrors.New("invalid")
	ErrNotBoolean      = oerrors.New("not boolean")
	ErrNotComparable   = oerrors.New("not comparable")
	ErrNotDefined      = oerrors.New("not defined")
	ErrNotImplemented  = oerrors.New("not implemented")
	ErrNotMapping      = oerrors.New("not a mapping")
	ErrSyntaxError     = oerrors.New("syntax error")
	ErrUnexpectedValue = oerrors.New("unexpected value")
	ErrUnsupported     = oerrors.New("unsupported")
)

func (n *node) locatedError(err error) error {
	return errors.Located(n.path, n.Node, err)
}

func (m mappingChild) locatedError(err error) error {
	return errors.Located(m.path, m.value, err)
}
