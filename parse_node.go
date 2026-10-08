package yaml2yaml

import (
	"encoding/base64"
	"errors"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/directives"
	yerrors "github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

func parseNode(basePath string, src *yaml.Node) (node.Source, error) {
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
		return nil, yerrors.ErrEmpty
	}
}

func parseWithoutTag(basePath string, src *yaml.Node) (node.Source, error) {
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

func parseSequence(basePath string, src *yaml.Node) (node.SequenceSource, error) {
	result := node.SequenceSource(make([]node.Source, 0, len(src.Content)))

	for _, child := range src.Content {
		childNode, err := parseNode(basePath, child)
		if err != nil {
			continue
		}

		result = append(result, childNode)
	}

	if len(result) == 0 {
		return nil, yerrors.ErrEmpty
	}

	return result, nil
}

func parseMapping(basePath string, src *yaml.Node) (node.MappingSource, error) {
	result := make(node.MappingSource, 0, len(src.Content)/2)

	var errs []error

	for i := 0; i < len(src.Content); i += 2 {
		keyNode := src.Content[i]

		value, err := parseNode(basePath, src.Content[i+1])
		if err != nil {
			if !errors.Is(err, yerrors.ErrEmpty) {
				// TODO Locate error
				errs = append(errs, err)
			}

			continue
		}

		entry, err := parseMappingEntry(basePath, keyNode, value)
		if err != nil {
			if !errors.Is(err, yerrors.ErrEmpty) {
				// TODO Locate error
				errs = append(errs, err)
			}

			continue
		}

		result = append(result, entry)
	}

	return result, errors.Join(errs...)
}

func parseMappingEntry(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	switch keyNode.Tag {
	case directives.NameApply:
		name, err := parseWithoutTag(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return directives.Apply(name, value), nil
	case directives.NameDefine:
		name, err := parseWithoutTag(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return directives.Define(name, value), nil
	case directives.NameFor:
		sequence, err := parseWithoutTag(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return directives.For(sequence, value), nil
	case directives.NameIf:
		condition, err := parseWithoutTag(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return directives.If(condition, value)
	case directives.NameInclude:
		path, err := parseWithoutTag(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return directives.Include(basePath, path, parseNode), nil
	default:
		key, err := parseNode(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return node.MappingEntrySource{Key: key, Value: value}, nil
	}
}

func parseScalar(basePath string, src *yaml.Node) (node.Source, error) {
	switch src.Tag {
	case directives.NameApply:
		name, err := parseScalarWithoutTag(src)
		if err != nil {
			return nil, err
		}

		return directives.Apply(name, nil), nil
	case directives.NameInclude:
		path, err := parseScalarWithoutTag(src)
		if err != nil {
			return nil, err
		}

		return directives.Include(basePath, path, parseNode), nil
	case "!!bool":
		if value, err := strconv.ParseBool(src.Value); err == nil {
			return node.Bool(value), nil
		}
	case "!!int":
		if value, err := strconv.ParseInt(src.Value, 10, 64); err == nil {
			return node.Int(value), nil
		}
	case "!!float":
		if value, err := strconv.ParseFloat(src.Value, 64); err == nil {
			return node.Float(value), nil
		}
	case "!!binary":
		if value, err := base64.StdEncoding.DecodeString(src.Value); err == nil {
			return node.Binary(value), nil
		}
	case "!!str":
		if src.Value != "" && src.Value[0] == '?' {
			return expressions.Parse(src.Value[1:])
		}

		return node.String(src.Value), nil
	}

	return parseScalarWithoutTag(src)
}

func parseScalarWithoutTag(src *yaml.Node) (node.Source, error) {
	if src.Value != "" && src.Value[0] == '?' {
		return expressions.Parse(src.Value[1:])
	}

	if value, err := strconv.ParseInt(src.Value, 10, 64); err == nil {
		return node.Int(value), nil
	}

	if value, err := strconv.ParseFloat(src.Value, 64); err == nil {
		return node.Float(value), nil
	}

	if value, err := strconv.ParseBool(src.Value); err == nil {
		return node.Bool(value), nil
	}

	return node.String(src.Value), nil
}
