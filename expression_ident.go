package yaml2yaml

import "go/ast"

func (n *node) resolveIdent(expr *ast.Ident) (any, error) {
	switch expr.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		def := n.getDefinition(expr.Name)
		if def == nil {
			//nolint:nilnil // Nil as a value is useful here.
			return nil, nil
		}

		node, pErr := n.sameLevel(def).process()
		value, vErr := nodeToAny(node)

		return value, joinErrors(pErr, vErr)
	}
}
