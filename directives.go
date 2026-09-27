package yaml2yaml

import "go.yaml.in/yaml/v4"

const (
	dirApply  = "!apply"
	dirDefine = "!define"
	dirFor    = "!for"
	dirIf     = "!if"
)

//nolint:gochecknoglobals // List of directives tags, never modified...
var directivesTags = []string{
	dirApply,
	dirDefine,
	dirFor,
	dirIf,
}

func (n *node) mappingDirective(key, value *node) ([]mappingChild, error) {
	switch key.Tag {
	case dirApply:
		return n.applyMapping(key, value)
	case dirDefine:
		return n.defineMapping(key, value)
	case dirFor:
		return n.forMapping(key, value)
	case dirIf:
		return n.ifMapping(key, value)
	default:
		keyResult, kErr := key.process()
		valueResult, vErr := value.process()

		return []mappingChild{{path: key.path, key: keyResult, value: valueResult}}, joinErrors(kErr, vErr)
	}
}

func (n *node) scalarDirective() (*yaml.Node, error) {
	switch n.Tag {
	case dirApply:
		return n.applyScalar()
	default:
		return n.scalarNoDirective()
	}
}
