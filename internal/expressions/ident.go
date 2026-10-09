package expressions

import (
	"go/ast"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) ident(defs node.Definitions, expr *ast.Ident) (node.Node, error) {
	switch expr.Name {
	case "true":
		return node.NewBool(s.pos, true), nil
	case "false":
		return node.NewBool(s.pos, false), nil
	default:
		def := defs.Get(expr.Name)
		if def == nil {
			return nil, errors.ErrEmpty
		}

		return def.Process(defs)
	}
}
