package yaml2yaml

import (
	"encoding/base64"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

func anyToString(v any) string {
	switch v := v.(type) {
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	case []any:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = anyToString(item)
		}
		return strings.Join(result, ", ")
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func anyToNode(src *yaml.Node, v any) *yaml.Node {
	switch v := v.(type) {
	case bool:
		return scalarNode(src, "!!bool", strconv.FormatBool(v))
	case int:
		return scalarNode(src, "!!int", strconv.Itoa(v))
	case int8:
		return scalarNode(src, "!!int", strconv.FormatInt(int64(v), 10))
	case int16:
		return scalarNode(src, "!!int", strconv.FormatInt(int64(v), 10))
	case int32:
		return scalarNode(src, "!!int", strconv.FormatInt(int64(v), 10))
	case int64:
		return scalarNode(src, "!!int", strconv.FormatInt(v, 10))
	case float32:
		return scalarNode(src, "!!float", strconv.FormatFloat(float64(v), 'f', -1, 64))
	case float64:
		return scalarNode(src, "!!float", strconv.FormatFloat(v, 'f', -1, 64))
	case string:
		return scalarNode(src, "!!str", v)
	case []byte:
		return scalarNode(src, "!!binary", base64.StdEncoding.EncodeToString(v))
	case []any:
		content := make([]*yaml.Node, 0, len(v))
		for _, item := range v {
			if child := anyToNode(src, item); child != nil {
				content = append(content, child)
			}
		}

		return &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: content,
		}
	case map[string]any:
		content := make([]*yaml.Node, 0, 2*len(v))
		for _, key := range slices.Sorted(maps.Keys(v)) {
			if child := anyToNode(src, v[key]); child != nil {
				content = append(content, scalarNode(src, "!!str", key), child)
			}
		}

		return &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: content,
		}
	default:
		return scalarNode(src, src.Tag, fmt.Sprintf("%v", v))
	}
}

func scalarNode(src *yaml.Node, tag, value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   tag,
		Value: value,

		Line:   src.Line,
		Column: src.Column,
	}
}

func nodeToAny(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		return nodeToAny(node.Content[0])

	case yaml.AliasNode:
		return nodeToAny(node.Alias)

	case yaml.SequenceNode:
		return sequenceToAny(node)

	case yaml.MappingNode:
		return mappingToAny(node)

	case yaml.ScalarNode:
		return scalarToAny(node)

	default:
		return nil, fmt.Errorf("%w kind: %v", ErrUnsupported, kindString(node.Kind))
	}
}

func sequenceToAny(node *yaml.Node) (any, error) {
	var errs []error

	children := make([]any, len(node.Content))

	for i, child := range node.Content {
		child, err := nodeToAny(child)
		if err != nil {
			errs = append(errs, err)
		}

		children[i] = child
	}

	return children, joinErrors(errs...)
}

func mappingToAny(node *yaml.Node) (any, error) {
	var errs []error

	children := make(map[string]any, len(node.Content)/2)

	for i := 0; i < len(node.Content); i += 2 {
		keyAny, err := scalarToAny(node.Content[i])
		if err != nil {
			errs = append(errs, err)
		}

		key := anyToString(keyAny)

		child, err := nodeToAny(node.Content[i+1])
		if err != nil {
			errs = append(errs, err)
		}

		children[key] = child
	}

	return children, joinErrors(errs...)
}

func scalarToAny(node *yaml.Node) (any, error) {
	switch node.Tag {
	case "!!bool":
		return strconv.ParseBool(node.Value)

	case "!!int":
		return strconv.Atoi(node.Value)

	case "!!float":
		return strconv.ParseFloat(node.Value, 64)

	case "!!str":
		return node.Value, nil

	case "!!binary":
		return base64.StdEncoding.DecodeString(node.Value)

	default:
		return node.Value, nil
	}
}
