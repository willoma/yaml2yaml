package expressions

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/internal/node"
)

func ident(defs node.Definitions, expr *ast.Ident) (node.Node, error) {
	switch expr.Name {
	case "true":
		return node.Bool(true), nil
	case "false":
		return node.Bool(false), nil
	default:
		def := defs.Get(expr.Name)
		if def == nil {
			//nolint:nilnil // Nil as a value is useful here.
			return nil, nil
		}

		return def.Process(defs)
	}
}
