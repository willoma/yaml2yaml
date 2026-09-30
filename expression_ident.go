package yaml2yaml

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
)

func resolveIdent(defs *definitions.Definitions, expr *ast.Ident) (any, error) {
	switch expr.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		def := defs.Get(expr.Name)
		if def == nil {
			//nolint:nilnil // Nil as a value is useful here.
			return nil, nil
		}

		node, err := processTransparent(defs.Child(), def)

		return convert.NodeToAny(node), err
	}
}
