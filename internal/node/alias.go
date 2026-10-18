package node

type Alias struct {
	pos    [2]int
	anchor string
	Target Source
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
		return a.anchor == otherAlias.anchor && a.Target.EqualsSource(otherAlias.Target)
	}

	return false
}

func (a *Alias) Process(defs Definitions) (Node, error) {
	return a.Target.Process(defs)
}
