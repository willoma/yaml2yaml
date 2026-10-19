package directives

import (
	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

const NameDefine = "!define"

type dirDefine struct {
	name  node.Source
	value node.Source
}

func Define(name, value node.Source) *dirDefine {
	return &dirDefine{name: name, value: value}
}

func (d *dirDefine) Pos() [2]int {
	return d.name.Pos()
}

func (d *dirDefine) String() string {
	return NameDefine + " " + d.name.String()
}

func (d *dirDefine) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirDefine); ok {
		return d.name.EqualsSource(other.name) && d.value.EqualsSource(other.value)
	}

	return false
}

func (d *dirDefine) Process(defs node.Definitions) (node.Node, error) {
	nameNode, err := d.name.Process(defs)
	if err != nil {
		return nil, err
	}

	defs.Set(nameNode.String(), d.value)

	return nil, errors.ErrEmpty
}

func (d *dirDefine) PopulateDefs(node.Definitions) error {
	return nil
}
