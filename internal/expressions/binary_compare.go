package expressions

import (
	"go/token"

	"github.com/willoma/yaml2yaml/internal/node"
)

func binCompare(left, right node.Node, op token.Token) bool {
	// Before trying scalars, try sequences, mappings, and binary.
	if lseq, ok := left.(node.Sequence); ok {
		if rseq, ok := right.(node.Sequence); ok {
			return compareSequences(lseq, rseq, op)
		}
	}

	if lmap, ok := left.(node.Mapping); ok {
		if rmap, ok := right.(node.Mapping); ok {
			return compareMappings(lmap, rmap, op)
		}
	}

	if lbin, ok := left.(node.Binary); ok {
		if rbin, ok := right.(node.Binary); ok {
			return compareBinaries(lbin, rbin, op)
		}
	}

	// First try int, which is the most restrictive.
	if lint, ok := left.Int(); ok {
		if rint, ok := right.Int(); ok {
			return compareTyped(lint, rint, op)
		}
	}

	// Next try float, which is less restrictive than int.
	if lfloat, ok := left.Float(); ok {
		if rfloat, ok := right.Float(); ok {
			return compareTyped(lfloat, rfloat, op)
		}
	}

	// Next try bool, which is less restrictive than float.
	if lbool, ok := left.Bool(); ok {
		if rbool, ok := right.Bool(); ok {
			return compareBools(lbool, rbool, op)
		}
	}

	// Finally, no restrictive type matched, compare as strings...
	return compareTyped(left.String(), right.String(), op)
}

func compareSequences(left, right node.Sequence, op token.Token) bool {
	switch {
	case len(left) > len(right):
		return op == token.NEQ || op == token.GTR || op == token.GEQ
	case len(left) < len(right):
		return op == token.NEQ || op == token.LSS || op == token.LEQ
	}

	for i := range left {
		if !binCompare(left[i], right[i], op) {
			return false
		}
	}

	return true
}

func compareMappings(left, right node.Mapping, op token.Token) bool {
	switch {
	case len(left) > len(right):
		return op == token.NEQ || op == token.GTR || op == token.GEQ
	case len(left) < len(right):
		return op == token.NEQ || op == token.LSS || op == token.LEQ
	}

	for i := range left {
		if !binCompare(left[i].Key(), right[i].Key(), op) {
			return false
		}

		if !binCompare(left[i].Value(), right[i].Value(), op) {
			return false
		}
	}

	return true
}

func compareBools(left, right bool, op token.Token) bool {
	switch op {
	case token.EQL, token.LEQ, token.GEQ:
		return left == right
	case token.NEQ:
		return left != right
	default:
		return false
	}
}

func compareBinaries(left, right node.Binary, op token.Token) bool {
	return compareTyped(left.Compare(right), 0, op)
}

func compareTyped[T ~int | ~int64 | ~float64 | ~string](left, right T, op token.Token) bool {
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
