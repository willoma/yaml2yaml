package expressions

import (
	"fmt"
	"go/token"
	"slices"
	"strings"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) binMath(left, right node.Node, op token.Token) (node.Node, error) {
	if lseq, ok := left.(node.Sequence); ok {
		return s.mathSequence(lseq, right, op)
	}

	if lmap, ok := left.(node.Mapping); ok {
		if rmap, ok := right.(node.Mapping); ok {
			return s.mathMappings(lmap, rmap, op)
		}
	}

	if lbin, ok := left.(node.Binary); ok {
		return s.mathBinary(lbin, right, op)
	}

	if lint, ok := left.Int(); ok {
		if rint, ok := right.Int(); ok {
			return s.mathInts(lint, rint, op)
		}
	}

	if lfloat, ok := left.Float(); ok {
		if rfloat, ok := right.Float(); ok {
			return s.mathFloats(lfloat, rfloat, op)
		}
	}

	return s.mathString(left.String(), right, op)
}

func (s *Source) mathSequence(left node.Sequence, right node.Node, op token.Token) (node.Node, error) {
	switch op {
	case token.ADD:
		if rseq, ok := right.(node.Sequence); ok {
			left = append(left, rseq...)
		} else {
			left = append(left, right)
		}

		return left, nil
	case token.MUL:
		rint, ok := right.Int()
		if !ok {
			return nil, fmt.Errorf("right: %w", errors.ErrNotANumber)
		}

		return slices.Repeat(left, int(rint)), nil

	default:
		return nil, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithSequences)
	}
}

func (s *Source) mathBinary(left node.Binary, right node.Node, op token.Token) (node.Binary, error) {
	switch op {
	case token.ADD:
		rbin, ok := right.(node.Binary)
		if !ok {
			return node.Binary{}, fmt.Errorf("right: %w", errors.ErrNotABinary)
		}

		return left.Append(rbin), nil
	case token.MUL:
		rint, ok := right.Int()
		if !ok {
			return node.Binary{}, fmt.Errorf("right: %w", errors.ErrNotANumber)
		}

		return left.Repeat(int(rint)), nil
	default:
		return node.Binary{}, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithBinaries)
	}
}

func (s *Source) mathMappings(left, right node.Mapping, op token.Token) (node.Node, error) {
	if op == token.ADD {
		return left.Merge(right)
	}

	return nil, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithMaps)
}

func (s *Source) mathInts(left, right int64, op token.Token) (node.Int, error) {
	switch op {
	case token.ADD:
		return node.NewInt(s.pos, left+right), nil
	case token.SUB:
		return node.NewInt(s.pos, left-right), nil
	case token.MUL:
		return node.NewInt(s.pos, left*right), nil
	case token.QUO:
		return node.NewInt(s.pos, left/right), nil
	case token.REM:
		return node.NewInt(s.pos, left%right), nil
	default:
		return node.Int{}, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithIntegers)
	}
}

func (s *Source) mathFloats(left, right float64, op token.Token) (node.Float, error) {
	switch op {
	case token.ADD:
		return node.NewFloat(s.pos, left+right), nil
	case token.SUB:
		return node.NewFloat(s.pos, left-right), nil
	case token.MUL:
		return node.NewFloat(s.pos, left*right), nil
	case token.QUO:
		return node.NewFloat(s.pos, left/right), nil
	default:
		return node.Float{}, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithFloats)
	}
}

func (s *Source) mathString(left string, right node.Node, op token.Token) (node.String, error) {
	switch op {
	case token.ADD:
		return node.NewString(s.pos, left+right.String()), nil
	case token.MUL:
		rint, ok := right.Int()
		if !ok {
			return node.String{}, fmt.Errorf("right: %w", errors.ErrNotANumber)
		}

		return node.NewString(s.pos, strings.Repeat(left, int(rint))), nil
	default:
		return node.String{}, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithStrings)
	}
}
