package expressions

import (
	"go/token"

	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) logical(left, right node.Node, op token.Token) node.Bool {
	bleft, _ := left.Bool()
	bright, _ := right.Bool()

	if op == token.LAND {
		return node.NewBool(s.pos, bleft && bright)
	}

	return node.NewBool(s.pos, bleft || bright)
}
