package expressions

import (
	"go/token"

	"github.com/willoma/yaml2yaml/internal/node"
)

func logical(left, right node.Node, op token.Token) bool {
	bleft, _ := left.Bool()
	bright, _ := right.Bool()

	if op == token.LAND {
		return bleft && bright
	}

	return bleft || bright
}
