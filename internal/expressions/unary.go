package expressions

// import (
// 	"fmt"
// 	"go/ast"
// 	"go/token"

// 	"github.com/willoma/yaml2yaml/internal/errors"
// 	"github.com/willoma/yaml2yaml/internal/node"
// )

// func unary(defs *node.Definitions, expr *ast.UnaryExpr) (node.Node, error) {
// 	value, rErr := Resolve(defs, expr.X)

// 	switch expr.Op {
// 	case token.NOT:
// 		ok, _ := value.Bool()

// 		return node.Bool(!ok), rErr
// 	default:
// 		tErr := fmt.Errorf("%q: %w in unary expression", expr.Op.String(), errors.ErrUnsupportedOperator)

// 		return nil, errors.Join(rErr, tErr)
// 	}
// }
