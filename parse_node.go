package yaml2yaml

import (
	"encoding/base64"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/directives"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

func parseNode(basePath string, src *yaml.Node) node.Source {
	switch src.Kind {
	case yaml.DocumentNode:
		return parseNode(basePath, src.Content[0])
	case yaml.AliasNode:
		return parseNode(basePath, src.Alias)
	case yaml.SequenceNode:
		return parseSequence(basePath, src)
	case yaml.MappingNode:
		return parseMapping(basePath, src)
	case yaml.ScalarNode:
		return parseScalar(basePath, src)
	default:
		return nil
	}
}

func parseWithoutTag(basePath string, src *yaml.Node) node.Source {
	return parseNode(basePath, &yaml.Node{
		Kind:    src.Kind,
		Value:   src.Value,
		Anchor:  src.Anchor,
		Alias:   src.Alias,
		Content: src.Content,
		Line:    src.Line,
		Column:  src.Column,
	})
}

func parseSequence(basePath string, src *yaml.Node) node.SequenceSource {
	result := node.SequenceSource(make([]node.Source, len(src.Content)))

	for i, child := range src.Content {
		result[i] = parseNode(basePath, child)
	}

	return result
}

func parseMapping(path string, src *yaml.Node) node.MappingSource {
	result := make(node.MappingSource, len(src.Content)/2)

	for i := 0; i < len(src.Content); i += 2 {
		keyNode := src.Content[i]
		value := parseNode(path, src.Content[i+1])

		switch keyNode.Tag {
		case directives.NameApply:
			result[i/2] = &directives.Apply{Name: parseWithoutTag(path, keyNode), Args: value}
		case directives.NameDefine:
			result[i/2] = &directives.Define{Name: parseWithoutTag(path, keyNode), Value: value}
		case directives.NameFor:
			result[i/2] = &directives.For{Sequence: parseWithoutTag(path, keyNode), Value: value}
		case directives.NameIf:
			result[i/2] = &directives.If{Condition: parseWithoutTag(path, keyNode), Value: value}
		case directives.NameInclude:
			result[i/2] = &directives.Include{BasePath: path, Path: parseWithoutTag(path, keyNode), ParseNode: parseNode}
		default:
			result[i/2] = node.MappingEntrySource{Key: parseNode(path, keyNode), Value: value}
		}
	}

	return result
}

func parseScalar(basePath string, src *yaml.Node) node.Source {
	switch src.Tag {
	case directives.NameApply:
		return &directives.Apply{Name: parseScalarWithoutTag(src)}
	case directives.NameInclude:
		return &directives.Include{BasePath: basePath, Path: parseScalarWithoutTag(src), ParseNode: parseNode}
	case "!!bool":
		if value, err := strconv.ParseBool(src.Value); err == nil {
			return node.Bool(value)
		}
	case "!!int":
		if value, err := strconv.ParseInt(src.Value, 10, 64); err == nil {
			return node.Int(value)
		}
	case "!!float":
		if value, err := strconv.ParseFloat(src.Value, 64); err == nil {
			return node.Float(value)
		}
	case "!!binary":
		if value, err := base64.StdEncoding.DecodeString(src.Value); err == nil {
			return node.Binary(value)
		}
	case "!!str":
		if src.Value != "" && src.Value[0] == '?' {
			return expressions.Source(src.Value[1:])
		}

		return node.String(src.Value)
	}

	return parseScalarWithoutTag(src)
}

func parseScalarWithoutTag(src *yaml.Node) node.Source {
	if src.Value != "" && src.Value[0] == '?' {
		return expressions.Source(src.Value[1:])
	}

	if value, err := strconv.ParseInt(src.Value, 10, 64); err == nil {
		return node.Int(value)
	}

	if value, err := strconv.ParseFloat(src.Value, 64); err == nil {
		return node.Float(value)
	}

	if value, err := strconv.ParseBool(src.Value); err == nil {
		return node.Bool(value)
	}

	return node.String(src.Value)
}
