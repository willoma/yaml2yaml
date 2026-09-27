package yaml2yaml

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
)

func (n *node) resolveIdent(expr *ast.Ident) (any, error) {
	switch expr.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		def := n.definitions.Get(expr.Name)
		if def == nil {
			//nolint:nilnil // Nil as a value is useful here.
			return nil, nil
		}

		node, pErr := n.sameLevel(def).process()
		value, vErr := nodeToAny(node)

		return value, errors.Join(pErr, vErr)
	}
}
