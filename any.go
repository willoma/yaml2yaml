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

func anyToString(value any) string {
	switch value := value.(type) {
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case string:
		return value
	case []any:
		result := make([]string, len(value))
		for i, item := range value {
			result[i] = anyToString(item)
		}

		return strings.Join(result, ", ")
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", value)
	}
}

func anyToNode(src *yaml.Node, value any) *yaml.Node {
	switch value := value.(type) {
	case bool:
		return scalarNode(src, "!!bool", strconv.FormatBool(value))
	case int:
		return scalarNode(src, "!!int", strconv.Itoa(value))
	case int8:
		return scalarNode(src, "!!int", strconv.FormatInt(int64(value), 10))
	case int16:
		return scalarNode(src, "!!int", strconv.FormatInt(int64(value), 10))
	case int32:
		return scalarNode(src, "!!int", strconv.FormatInt(int64(value), 10))
	case int64:
		return scalarNode(src, "!!int", strconv.FormatInt(value, 10))
	case float32:
		return scalarNode(src, "!!float", strconv.FormatFloat(float64(value), 'f', -1, 64))
	case float64:
		return scalarNode(src, "!!float", strconv.FormatFloat(value, 'f', -1, 64))
	case string:
		return scalarNode(src, "!!str", value)
	case []byte:
		return scalarNode(src, "!!binary", base64.StdEncoding.EncodeToString(value))
	case []any:
		content := make([]*yaml.Node, 0, len(value))
		for _, item := range value {
			if child := anyToNode(src, item); child != nil {
				content = append(content, child)
			}
		}

		return &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: content,
		}
	case map[string]any:
		content := make([]*yaml.Node, 0, 2*len(value))
		for _, key := range slices.Sorted(maps.Keys(value)) {
			if child := anyToNode(src, value[key]); child != nil {
				content = append(content, scalarNode(src, "!!str", key), child)
			}
		}

		return &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: content,
		}
	default:
		return scalarNode(src, src.Tag, fmt.Sprintf("%v", value))
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
	case yaml.StreamNode:
		return nil, fmt.Errorf("%w kind stream", ErrUnsupported)
	default:
		return nil, fmt.Errorf("%w unknown kind", ErrUnsupported)
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

//nolint:wrapcheck // Original errors are ok.
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
