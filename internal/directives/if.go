package directives

import (
	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameIf = "!if"

type If struct {
	Condition node.Source
	Value     node.Source
}

func (i *If) Process(defs node.Definitions) (node.Node, error) {
	condSource, cErr := i.Condition.Process(defs)

	condStr := condSource.String()

	exprNode := expressions.Source(condStr)

	cond, eErr := exprNode.Process(defs)

	if ok, _ := cond.Bool(); !ok {
		return nil, errors.Join(cErr, eErr)
	}

	result, pErr := i.Value.Process(defs)

	return result, errors.Join(cErr, eErr, pErr)
}
