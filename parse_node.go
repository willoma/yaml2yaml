package yaml2yaml

import (
	"encoding/base64"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/directives"
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
		return nil, errors.ErrEmpty
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

func parseSequence(basePath string, src *yaml.Node) (*node.SequenceSource, error) {
	var errs errors.JoinedError

	result := node.NewSequenceSource(pos(src), len(src.Content))

	for i, child := range src.Content {
		childNode, err := parseNode(basePath, child)
		if err != nil {
			errs = errs.Append(errors.LocateByIndex(i, pos(src), err))

			continue
		}

		result.Value = append(result.Value, childNode)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	if len(result.Value) == 0 {
		return nil, errors.ErrEmpty
	}

	return result, nil
}

func parseMapping(basePath string, src *yaml.Node) (*node.MappingSource, error) {
	result := node.NewMappingSource(pos(src), len(src.Content)/2)

	var errs errors.JoinedError

	for i := 0; i < len(src.Content); i += 2 {
		keyNode := src.Content[i]

		value, err := parseNode(basePath, src.Content[i+1])
		if err != nil {
			errs = errs.Append(locateMappingEntryError(i, basePath, keyNode, err))

			continue
		}

		entry, err := parseMappingEntry(basePath, keyNode, value)
		if err != nil {
			errs = errs.Append(locateMappingEntryError(i, basePath, keyNode, err))

			continue
		}

		result.Value = append(result.Value, entry)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}

func parseMappingEntry(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	switch keyNode.Tag {
	case directives.NameApply:
		return parseMappingDirectiveApply(basePath, keyNode, value)
	case directives.NameDefine:
		return parseMappingDirectiveDefine(basePath, keyNode, value)
	case directives.NameExpr:
		return parseMappingDirectiveExpr(basePath, keyNode, value)
	case directives.NameFmt:
		return parseMappingDirectiveFmt(basePath, keyNode, value)
	case directives.NameFor:
		return parseMappingDirectiveFor(basePath, keyNode, value)
	case directives.NameIf:
		return parseMappingDirectiveIf(basePath, keyNode, value)
	case directives.NameInclude:
		return parseMappingDirectiveInclude(basePath, keyNode)
	default:
		key, err := parseNode(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		return node.NewMappingEntrySource(key, value), nil
	}
}

func parseMappingDirectiveApply(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	name, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	return directives.Apply(name, value), nil
}

func parseMappingDirectiveDefine(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	name, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	return directives.Define(name, value), nil
}

func parseMappingDirectiveExpr(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	key, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	expr, err := directives.Expr(key)
	if err != nil {
		return nil, err
	}

	return node.NewMappingEntrySource(expr, value), nil
}

func parseMappingDirectiveFmt(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	key, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	fmtSrc, err := directives.Fmt(key)
	if err != nil {
		return nil, err
	}

	return node.NewMappingEntrySource(fmtSrc, value), nil
}

func parseMappingDirectiveFor(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	sequence, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	return directives.For(sequence, value), nil
}

func parseMappingDirectiveIf(basePath string, keyNode *yaml.Node, value node.Source) (node.Source, error) {
	condition, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	return directives.If(condition, value)
}

func parseMappingDirectiveInclude(basePath string, keyNode *yaml.Node) (node.Source, error) {
	path, err := parseWithoutTag(basePath, keyNode)
	if err != nil {
		return nil, err
	}

	return directives.Include(basePath, path, parseNode), nil
}

func parseScalar(basePath string, src *yaml.Node) (node.Source, error) {
	switch src.Tag {
	case directives.NameApply:
		name, err := parseScalarWithoutTag(src)
		if err != nil {
			return nil, err
		}

		return directives.Apply(name, nil), nil
	case directives.NameExpr:
		expr, err := parseScalarWithoutTag(src)
		if err != nil {
			return nil, err
		}

		return directives.Expr(expr)
	case directives.NameFmt:
		str, err := parseScalarWithoutTag(src)
		if err != nil {
			return nil, err
		}

		return directives.Fmt(str)
	case directives.NameInclude:
		path, err := parseScalarWithoutTag(src)
		if err != nil {
			return nil, err
		}

		return directives.Include(basePath, path, parseNode), nil
	case "!!bool":
		if value, err := strconv.ParseBool(src.Value); err == nil {
			return node.NewBool(pos(src), value), nil
		}
	case "!!int":
		if value, err := strconv.ParseInt(src.Value, 10, 64); err == nil {
			return node.NewInt(pos(src), value), nil
		}
	case "!!float":
		if value, err := strconv.ParseFloat(src.Value, 64); err == nil {
			return node.NewFloat(pos(src), value), nil
		}
	case "!!binary":
		if value, err := base64.StdEncoding.DecodeString(src.Value); err == nil {
			return node.NewBinary(pos(src), value), nil
		}
	case "!!str":
		return node.NewString(pos(src), src.Value), nil
	}

	return parseScalarWithoutTag(src)
}

func parseScalarWithoutTag(src *yaml.Node) (node.Source, error) {
	if value, err := strconv.ParseInt(src.Value, 10, 64); err == nil {
		return node.NewInt(pos(src), value), nil
	}

	if value, err := strconv.ParseFloat(src.Value, 64); err == nil {
		return node.NewFloat(pos(src), value), nil
	}

	if value, err := strconv.ParseBool(src.Value); err == nil {
		return node.NewBool(pos(src), value), nil
	}

	return node.NewString(pos(src), src.Value), nil
}
