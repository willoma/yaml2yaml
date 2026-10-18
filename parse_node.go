package yaml2yaml

import (
	"encoding/base64"
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/directives"
	"github.com/willoma/yaml2yaml/internal/node"
)

const mergeKey = "<<"

func parseNode(basePath string, src *yaml.Node) (node.Source, error) {
	switch src.Kind {
	case yaml.DocumentNode:
		return parseNode(basePath, src.Content[0])
	case yaml.AliasNode:
		return parseAlias(basePath, src)
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

func parseAlias(basePath string, src *yaml.Node) (*node.Alias, error) {
	target, err := parseNode(basePath, src.Alias)
	if err != nil {
		return nil, errors.LocateByKey(src.Anchor, pos(src), err)
	}

	return node.NewAlias(pos(src), src.Anchor, target), nil
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

		dirEntry, isDir, err := parseMappingDirective(basePath, keyNode, value)
		if err != nil {
			errs = errs.Append(locateMappingEntryError(i, basePath, keyNode, err))

			continue
		}

		if isDir {
			appendToMapping(result, dirEntry, true)

			continue
		}

		key, err := parseNode(basePath, keyNode)
		if err != nil {
			return nil, err
		}

		entry := node.NewMappingEntrySource(key, value)
		if entry.Key.String() == mergeKey {
			if wasAlias := appendAliasesToMapping(result, entry.Value); wasAlias {
				continue
			}
		}

		appendToMapping(result, entry, true)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}

func appendAliasesToMapping(result *node.MappingSource, value node.Source) bool {
	switch value := value.(type) {
	case *node.Alias:
		aliasMap, ok := value.Target.(*node.MappingSource)
		if !ok {
			return false
		}

		for _, item := range aliasMap.Value {
			appendToMapping(result, item, false)
		}

		return true

	case *node.SequenceSource:
		var items []node.Source

		for _, entry := range value.Value {
			alias, ok := entry.(*node.Alias)
			if !ok {
				return false
			}

			aliasMap, ok := alias.Target.(*node.MappingSource)
			if !ok {
				return false
			}

			items = append(items, aliasMap.Value...)
		}

		for _, item := range items {
			appendToMapping(result, item, false)
		}

		return true

	default:
		return false
	}
}

func appendToMapping(result *node.MappingSource, item node.Source, overwrite bool) {
	if newEntry, ok := item.(*node.MappingEntrySource); ok {
		for _, item := range result.Value {
			entry, ok := item.(*node.MappingEntrySource)
			if !ok {
				continue
			}

			if entry.Key.EqualsSource(newEntry.Key) {
				if overwrite {
					entry.Value = newEntry.Value
				}

				return
			}
		}
	}

	result.Value = append(result.Value, item)
}

func parseMappingDirective(
	basePath string, keyNode *yaml.Node, value node.Source,
) (result node.Source, isDir bool, err error) {
	switch keyNode.Tag {
	case directives.NameApply:
		result, err = parseMappingDirectiveApply(basePath, keyNode, value)
	case directives.NameDefine:
		result, err = parseMappingDirectiveDefine(basePath, keyNode, value)
	case directives.NameExpr:
		result, err = parseMappingDirectiveExpr(basePath, keyNode, value)
	case directives.NameFmt:
		result, err = parseMappingDirectiveFmt(basePath, keyNode, value)
	case directives.NameFor:
		result, err = parseMappingDirectiveFor(basePath, keyNode, value)
	case directives.NameIf:
		result, err = parseMappingDirectiveIf(basePath, keyNode, value)
	case directives.NameInclude:
		result, err = parseMappingDirectiveInclude(basePath, keyNode)
	}

	if err != nil {
		return nil, false, err
	}

	return result, result != nil, nil
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
