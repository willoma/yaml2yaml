package convert

import (
	"encoding/base64"
	"fmt"
	"maps"
	"slices"
	"strconv"

	"go.yaml.in/yaml/v4"
)

// AnyToNode converts an any value to a [*yaml.Node].
func AnyToNode(value any) *yaml.Node {
	switch value := value.(type) {
	case bool:
		return scalarNode("!!bool", strconv.FormatBool(value))
	case int:
		return scalarNode("!!int", strconv.Itoa(value))
	case int8:
		return scalarNode("!!int", strconv.FormatInt(int64(value), 10))
	case int16:
		return scalarNode("!!int", strconv.FormatInt(int64(value), 10))
	case int32:
		return scalarNode("!!int", strconv.FormatInt(int64(value), 10))
	case int64:
		return scalarNode("!!int", strconv.FormatInt(value, 10))
	case float32:
		return scalarNode("!!float", strconv.FormatFloat(float64(value), 'f', -1, 64))
	case float64:
		return scalarNode("!!float", strconv.FormatFloat(value, 'f', -1, 64))
	case string:
		return scalarNode("!!str", value)
	case []byte:
		return scalarNode("!!binary", base64.StdEncoding.EncodeToString(value))
	case []any:
		return sequenceNode(value)
	case map[string]any:
		return mappingNode(value)
	default:
		return scalarNode("", fmt.Sprintf("%v", value))
	}
}

func sequenceNode(value []any) *yaml.Node {
	content := make([]*yaml.Node, len(value))
	for i, item := range value {
		content[i] = AnyToNode(item)
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Content: content}
}

func mappingNode(value map[string]any) *yaml.Node {
	content := make([]*yaml.Node, 0, 2*len(value))
	for _, key := range slices.Sorted(maps.Keys(value)) {
		content = append(content, scalarNode("!!str", key), AnyToNode(value[key]))
	}

	return &yaml.Node{Kind: yaml.MappingNode, Content: content}
}

func scalarNode(tag, value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: value}
}
