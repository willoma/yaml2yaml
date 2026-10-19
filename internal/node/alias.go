package node

import "github.com/willoma/yaml2yaml/errors"

type Alias struct {
	pos    [2]int
	anchor string
	target Source
}

func NewAlias(pos [2]int, anchor string, target Source) *Alias {
	return &Alias{pos, anchor, target}
}

func (a *Alias) Pos() [2]int {
	return a.pos
}

func (a *Alias) String() string {
	return a.anchor
}

func (a *Alias) EqualsSource(other Source) bool {
	if otherAlias, ok := other.(*Alias); ok {
		return a.anchor == otherAlias.anchor && a.target.EqualsSource(otherAlias.target)
	}

	return false
}

func (a *Alias) Process(defs Definitions) (Node, error) {
	return a.target.Process(defs)
}

func (a *Alias) SourceMap(defs Definitions) (map[string]Source, error) {
	sourceMapper, ok := a.target.(SourceMapper)
	if !ok {
		return nil, errors.ErrNotAMapping
	}

	return sourceMapper.SourceMap(defs)
}
