package yaml2yaml

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

type node struct {
	*yaml.Node

	path        string
	parent      *node
	definitions map[string]*yaml.Node
}

func newNode(src *yaml.Node) *node {
	return &node{
		Node:        src,
		path:        "",
		parent:      nil,
		definitions: make(map[string]*yaml.Node),
	}
}

func (n *node) child(pathSuffix string, src *yaml.Node) *node {
	return &node{
		Node:        src,
		path:        n.path + pathSuffix,
		parent:      n,
		definitions: make(map[string]*yaml.Node),
	}
}

func (n *node) indexedChild(i int, src *yaml.Node) *node {
	return n.child("["+strconv.Itoa(i)+"]", src)
}

func (n *node) keyedChild(tag, key string, src *yaml.Node) *node {
	if slices.Contains(directivesTags, tag) {
		key = tag + " " + key
	}

	pathSuffix := key
	if n.path != "" {
		pathSuffix = "." + key
	}

	return n.child(pathSuffix, src)
}

func (n *node) process() (*yaml.Node, error) {
	switch n.Kind {
	case yaml.DocumentNode:
		return n.child("", n.Content[0]).process()
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
