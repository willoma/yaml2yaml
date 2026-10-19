package node

import "go.yaml.in/yaml/v4"

type (
	// Source is a parsed node source, before any yaml2yaml processing.
	Source interface {
		// Pos returns the position of the [Source] in the source file.
		Pos() [2]int
		// String returns the [Source] representation as a string.
		String() string

		// EqualsSource compares the [Source] to another [Source].
		EqualsSource(other Source) bool
		// Process generates a [Node] from the [Source].
		Process(defs Definitions) (Node, error)
	}

	// Node is a processed node, ready to be manipulated and converted back to a [yaml.Node].
	Node interface {
		// Bool returns the node value as a boolean, and false if the value was not parseable.
		Bool() (value, ok bool)
		// Int returns the node value as an integer, and false if the value was not exactly parseable.
		Int() (int64, bool)
		// Float returns the node value as a float, and false if the value was not parseable.
		Float() (float64, bool)
		// String returns the node value as a string.
		String() string

		// EqualsNode compares the [Node] to another [Node].
		EqualsNode(other Node) bool
		// ToYAML generates a [yaml.Node] from the [Node].
		ToYAML() *yaml.Node
	}

	SourceMapper interface {
		SourceMap(defs Definitions) (map[string]Source, error)
	}
)
