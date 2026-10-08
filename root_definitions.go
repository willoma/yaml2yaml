package yaml2yaml

import "github.com/willoma/yaml2yaml/internal/node"

type rootDefinitions map[string]node.Source

func (d rootDefinitions) Get(name string) node.Source { return d[name] }

func (d rootDefinitions) Set(name string, value node.Source) { d[name] = value }
