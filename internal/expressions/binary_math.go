package expressions

import (
	"fmt"
	"go/token"

	"github.com/willoma/yaml2yaml/internal/errors"
	"github.com/willoma/yaml2yaml/internal/node"
)

func (s *Source) binMath(left, right node.Node, op token.Token) (node.Node, error) {
	if lseq, ok := left.(node.Sequence); ok {
		if rseq, ok := right.(node.Sequence); ok {
			return s.mathSequences(lseq, rseq, op)
		}
	}

	if lmap, ok := left.(node.Mapping); ok {
		if rmap, ok := right.(node.Mapping); ok {
			return s.mathMappings(lmap, rmap, op)
		}
	}

	if lbin, ok := left.(node.Binary); ok {
		if rbin, ok := right.(node.Binary); ok {
			return s.mathBinaries(lbin, rbin, op)
		}
	}

	if lint, ok := left.Int(); ok {
		if rint, ok := right.Int(); ok {
			return s.mathInts(lint, rint, op), nil
		}
	}

	if lfloat, ok := left.Float(); ok {
		if rfloat, ok := right.Float(); ok {
			result, err := s.mathFloats(lfloat, rfloat, op)

			return node.NewFloat(s.pos, result), err
		}
	}

	result, err := s.mathStrings(left.String(), right.String(), op)

	return node.NewString(s.pos, result), err
}

func (s *Source) mathSequences(left, right node.Sequence, op token.Token) (node.Node, error) {
	if op == token.ADD {
		result := make(node.Sequence, 0, len(left)+len(right))
		result = append(result, left...)
		result = append(result, right...)

		return result, nil
	}

	return left, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithSequences)
}

func (s *Source) mathMappings(left, right node.Mapping, op token.Token) (node.Node, error) {
	if op == token.ADD {
		result, err := left.Merge(right)

		return result, err
	}

	return left, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithMaps)
}

func (s *Source) mathInts(left, right int64, op token.Token) node.Int {
	switch op {
	case token.ADD:
		return node.NewInt(s.pos, left+right)
	case token.SUB:
		return node.NewInt(s.pos, left-right)
	case token.MUL:
		return node.NewInt(s.pos, left*right)
	case token.QUO:
		return node.NewInt(s.pos, left/right)
	case token.REM:
		return node.NewInt(s.pos, left%right)
	default:
		return node.NewInt(s.pos, 0)
	}
}

func (s *Source) mathFloats(left, right float64, op token.Token) (float64, error) {
	switch op {
	case token.ADD:
		return left + right, nil
	case token.SUB:
		return left - right, nil
	case token.MUL:
		return left * right, nil
	case token.QUO:
		return left / right, nil
	default:
		return 0, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithFloats)
	}
}

func (s *Source) mathStrings(left, right string, op token.Token) (string, error) {
	switch op {
	case token.ADD:
		return left + right, nil
	default:
		return "", fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithStrings)
	}
}

func (s *Source) mathBinaries(left, right node.Binary, op token.Token) (node.Binary, error) {
	switch op {
	case token.ADD:
		return node.NewBinary(s.pos, append(left.Value, right.Value...)), nil
	default:
		return node.Binary{}, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithBinaries)
	}
}
