package yaml2yaml

import "go/ast"

func (n *node) resolveIdent(expr *ast.Ident) (any, error) {
	switch expr.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		def, dErr := n.getDefinition(expr.Name)
		if dErr != nil {
			return nil, dErr
		}

		node, pErr := n.child("", def).process()
		value, vErr := nodeToAny(node)

		return value, joinErrors(pErr, vErr)
	}
}
