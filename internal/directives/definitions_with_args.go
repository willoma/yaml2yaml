package directives

import (
	oerrors "errors"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

type definitionsWithArgs struct {
	parent      node.Definitions
	args        map[string]node.Source
	definitions map[string]node.Source
}

func defsWithArgs(parent node.Definitions, args node.Source) (*definitionsWithArgs, error) {
	defs := &definitionsWithArgs{
		parent:      parent,
		definitions: make(map[string]node.Source),
	}

	if args == nil {
		return defs, nil
	}

	sourceMapper, ok := args.(node.SourceMapper)
	if !ok {
		defs.args = map[string]node.Source{"_": args}

		return defs, nil
	}

	var err error

	defs.args, err = sourceMapper.SourceMap(defs)
	if oerrors.Is(err, errors.ErrEmpty) {
		defs.args = make(map[string]node.Source)
	} else if err != nil {
		return nil, err
	}

	defs.args["_"] = args

	return defs, nil
}

func (d *definitionsWithArgs) Get(name string) node.Source {
	if value, ok := d.args[name]; ok {
		return value
	}

	if value, ok := d.definitions[name]; ok {
		return value
	}

	return d.parent.Get(name)
}

func (d *definitionsWithArgs) Set(name string, value node.Source) {
	d.definitions[name] = value
}

func (d *definitionsWithArgs) passToParent() {
	for name, value := range d.definitions {
		d.parent.Set(name, value)
	}
}
