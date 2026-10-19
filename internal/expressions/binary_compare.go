package expressions

import (
	"go/token"

	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) binCompare(left, right node.Node, op token.Token) bool {
	// Before trying scalars, try sequences, mappings, and binary.
	if lseq, ok := left.(node.Sequence); ok {
		if rseq, ok := right.(node.Sequence); ok {
			return s.compareSequences(lseq, rseq, op)
		}
	}

	if lmap, ok := left.(node.Mapping); ok {
		if rmap, ok := right.(node.Mapping); ok {
			return s.compareMappings(lmap, rmap, op)
		}
	}

	if lbin, ok := left.(node.Binary); ok {
		if rbin, ok := right.(node.Binary); ok {
			return s.compareBinaries(lbin, rbin, op)
		}
	}

	// First try int, which is the most restrictive.
	if lint, ok := left.Int(); ok {
		if rint, ok := right.Int(); ok {
			return s.compareTyped(lint, rint, op)
		}
	}

	// Next try float, which is less restrictive than int.
	if lfloat, ok := left.Float(); ok {
		if rfloat, ok := right.Float(); ok {
			return s.compareTyped(lfloat, rfloat, op)
		}
	}

	// Next try bool, which is less restrictive than float.
	if lbool, ok := left.Bool(); ok {
		if rbool, ok := right.Bool(); ok {
			return s.compareBools(lbool, rbool, op)
		}
	}

	// Finally, no restrictive type matched, compare as strings...
	return s.compareTyped(left.String(), right.String(), op)
}

func (s *Source) compareSequences(left, right node.Sequence, op token.Token) bool {
	switch {
	case len(left) > len(right):
		return op == token.NEQ || op == token.GTR || op == token.GEQ
	case len(left) < len(right):
		return op == token.NEQ || op == token.LSS || op == token.LEQ
	}

	for i := range left {
		if !s.binCompare(left[i], right[i], op) {
			return false
		}
	}

	return true
}

func (s *Source) compareMappings(left, right node.Mapping, op token.Token) bool {
	switch {
	case len(left) > len(right):
		return op == token.NEQ || op == token.GTR || op == token.GEQ
	case len(left) < len(right):
		return op == token.NEQ || op == token.LSS || op == token.LEQ
	}

	for i := range left {
		if !s.binCompare(left[i].Key(), right[i].Key(), op) {
			return false
		}

		if !s.binCompare(left[i].Value(), right[i].Value(), op) {
			return false
		}
	}

	return true
}

func (s *Source) compareBools(left, right bool, op token.Token) bool {
	switch op {
	case token.EQL, token.LEQ, token.GEQ:
		return left == right
	case token.NEQ:
		return left != right
	default:
		return false
	}
}

func (s *Source) compareBinaries(left, right node.Binary, op token.Token) bool {
	result := left.Compare(right)

	switch op {
	case token.EQL:
		return result == 0
	case token.LSS:
		return result < 0
	case token.GTR:
		return result > 0
	case token.NEQ:
		return result != 0
	case token.LEQ:
		return result <= 0
	case token.GEQ:
		return result >= 0
	default:
		return false
	}
}

func (s *Source) compareTyped[T ~int64 | ~float64 | ~string](left, right T, op token.Token) bool {
	switch op {
	case token.EQL:
		return left == right
	case token.LSS:
		return left < right
	case token.GTR:
		return left > right
	case token.NEQ:
		return left != right
	case token.LEQ:
		return left <= right
	case token.GEQ:
		return left >= right
	default:
		return false
	}
}
