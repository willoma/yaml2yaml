package node

import "go.yaml.in/yaml/v4"

// Source is a parsed node source, before any yaml2yaml processing.
type Source interface {
	Pos() [2]int
	String() string

	EqualsSource(other Source) bool
	Process(defs Definitions) (Node, error)
}

// Node is a processed node, ready to be manipulated and converted back to a [yaml.Node].
type Node interface {
	// Bool returns the node value as a boolean, and false if the value was not parseable.
	Bool() (value, ok bool)
	// Int returns the node value as an integer, and false if the value was not exactly parseable.
	Int() (int64, bool)
	// Float returns the node value as a float, and false if the value was not parseable.
	Float() (float64, bool)
	// String returns the node value as a string.
	String() string

	EqualsNode(other Node) bool
	ToYAML() *yaml.Node
}
