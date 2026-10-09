package node

import "go.yaml.in/yaml/v4"

type Pos struct {
	Line   int
	Column int
}

func SourcePos(src *yaml.Node) Pos {
	return Pos{src.Line, src.Column}
}
