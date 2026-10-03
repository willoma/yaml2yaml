package directives

import (
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameDefine = "!define"

type Define struct {
	Name  node.Source
	Value node.Source
}

func (d *Define) Process(defs node.Definitions) (node.Node, error) {
	nameNode, nErr := d.Name.Process(defs)

	defs.Set(nameNode.String(), d.Value)

	return nil, nErr
}
