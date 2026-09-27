package yaml2yaml

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/definitions"
)

type node struct {
	*yaml.Node

	path        string
	definitions *definitions.Definitions
}

func newNode(src *yaml.Node) *node {
	return &node{
		Node:        src,
		path:        "",
		definitions: definitions.New(),
	}
}

func (n *node) sameLevel(src *yaml.Node) *node {
	return &node{
		Node:        src,
		path:        n.path,
		definitions: n.definitions.Sibling(),
	}
}

func (n *node) fileChild(path string, src *yaml.Node) *node {
	return &node{
		Node:        src,
		path:        path + ":" + n.path,
		definitions: n.definitions.Child(),
	}
}

func (n *node) indexedChild(i int, src *yaml.Node) *node {
	return &node{
		Node:        src,
		path:        n.path + "[" + strconv.Itoa(i) + "]",
		definitions: n.definitions.Child(),
	}
}

func (n *node) keyedChild(tag, key string, src *yaml.Node) *node {
	if slices.Contains(directivesTags, tag) {
		key = tag + " " + key
	}

	return &node{
		Node:        src,
		path:        childPath(n.path, key),
		definitions: n.definitions.Child(),
	}
}

func (n *node) process() (*yaml.Node, error) {
	switch n.Kind {
	case yaml.DocumentNode:
		return n.sameLevel(n.Content[0]).process()
	case yaml.AliasNode:
		return n.keyedChild("", "alias:"+n.Alias.Anchor, n.Alias).process()
	case yaml.SequenceNode:
		return n.processSequence()
	case yaml.MappingNode:
		return n.processMapping()
	case yaml.ScalarNode:
		return n.processScalar()
	case yaml.StreamNode:
		return nil, n.locatedError(fmt.Errorf("%w kind stream", ErrUnsupported))
	default:
		return nil, n.locatedError(fmt.Errorf("%w unnamed kind", ErrUnsupported))
	}
}

func (n *node) evaluateValue() (any, error) {
	if expr, ok := strings.CutPrefix(n.Value, "?"); ok {
		return n.resolveExpression(expr)
	}

	result, err := scalarToAny(n.Node)
	if err != nil {
		return nil, n.locatedError(err)
	}

	return result, nil
}

func childPath(path, key string) string {
	if path == "" {
		return key
	}

	return path + "." + key
}
