package yaml2yaml

import (
	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/convert"
	"github.com/willoma/yaml2yaml/definitions"
	"github.com/willoma/yaml2yaml/errors"
)

const (
	dirApply   = "!apply"
	dirDefine  = "!define"
	dirFor     = "!for"
	dirIf      = "!if"
	dirInclude = "!include"
)

func directiveMapping(
	defs *definitions.Definitions, key, value *yaml.Node,
) (result []mappingChild, processed bool, err error) {
	var resultNode *yaml.Node

	switch key.Tag {
	case dirApply:
		resultNode, err = doDirApply(defs, key, value)
	case dirDefine:
		err = doDirDefine(defs, key, value)
	case dirFor:
		resultNode, err = doDirFor(defs, key, value)
	case dirIf:
		resultNode, err = doDirIf(defs, key, value)
	case dirInclude:
		resultNode, err = doDirInclude(defs, key, value)
	default:
		return nil, false, nil
	}

	pathKey := key.Tag + " " + convert.NodeToString(key)
	err = errors.Locate(pathKey, key.Line, key.Column, err)

	switch {
	case resultNode == nil:
		result = nil
	case resultNode.Kind != yaml.MappingNode:
		result = []mappingChild{{pathKey: pathKey, value: resultNode}}
	default:
		result = make([]mappingChild, 0, len(resultNode.Content)/2)

		for i := 0; i < len(resultNode.Content); i += 2 {
			result = append(result, mappingChild{
				pathKey: convert.NodeToString(resultNode.Content[i]),
				key:     resultNode.Content[i],
				value:   resultNode.Content[i+1],
			})
		}
	}

	return result, true, err
}

func directiveScalar(defs *definitions.Definitions, node *yaml.Node) (result *yaml.Node, processed bool, err error) {
	switch node.Tag {
	case dirApply:
		result, err = doDirApply(defs, node, nil)
	case dirInclude:
		result, err = doDirInclude(defs, node, nil)
	default:
		return nil, false, nil
	}

	return result, true, err
}
