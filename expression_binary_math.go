package yaml2yaml

import (
	"bytes"
	"fmt"
	"go/token"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/willoma/yaml2yaml/errors"
)

func resolveMath(left, right any, op token.Token) (any, error) {
	switch left := left.(type) {
	case int:
		return resolveIntMath(left, right, op)
	case float64:
		return resolveFloatMath(left, right, op)
	case string:
		return resolveStringMath(left, right, op)
	case []byte:
		return resolveBinaryMath(left, right, op)
	case []any:
		return resolveSliceMath(left, right, op)
	case map[string]any:
		return resolveMapMath(left, right, op)
	default:
		return nil, fmt.Errorf("left (%T): %w", left, errors.ErrInvalidTypeForMathematicalOperation)
	}
}

func resolveIntMath(left int, right any, op token.Token) (any, error) {
	switch right := right.(type) {
	case int:
		return resolveIntsMath(left, right, op), nil
	case float64:
		return resolveFloatsMath(float64(left), right, op)
	case string:
		if iRight, err := strconv.Atoi(right); err == nil {
			return resolveIntsMath(left, iRight, op), nil
		}

		if fRight, err := strconv.ParseFloat(right, 64); err == nil {
			return resolveFloatsMath(float64(left), fRight, op)
		}

		if op == token.MUL {
			return strings.Repeat(right, left), nil
		}

		return resolveStringsMath(strconv.Itoa(left), right, op)
	case []byte:
		if op == token.MUL {
			return bytes.Repeat(right, left), nil
		}

		return resolveBinariesMath([]byte(strconv.Itoa(left)), right, op)
	case []any:
		if op == token.MUL {
			return slices.Repeat(right, left), nil
		}

		return resolveSlicesMath([]any{left}, right, op)

	default:
		return nil, fmt.Errorf("right (%T): %w", right, errors.ErrInvalidTypeForMathematicalOperation)
	}
}

func resolveIntsMath(left, right int, op token.Token) int {
	switch op {
	case token.ADD:
		return left + right
	case token.SUB:
		return left - right
	case token.MUL:
		return left * right
	case token.QUO:
		return left / right
	case token.REM:
		return left % right
	default:
		return 0
	}
}

func resolveFloatMath(left float64, right any, op token.Token) (any, error) {
	switch right := right.(type) {
	case int:
		return resolveFloatsMath(left, float64(right), op)
	case float64:
		return resolveFloatsMath(left, right, op)
	case string:
		if fRight, err := strconv.ParseFloat(right, 64); err == nil {
			return resolveFloatsMath(left, fRight, op)
		}

		return resolveStringsMath(strconv.FormatFloat(left, 'f', -1, 64), right, op)
	case []byte:
		return resolveBinariesMath([]byte(strconv.FormatFloat(left, 'f', -1, 64)), right, op)
	case []any:
		return resolveSlicesMath([]any{left}, right, op)
	default:
		return nil, fmt.Errorf("right (%T): %w", right, errors.ErrInvalidTypeForMathematicalOperation)
	}
}

func resolveFloatsMath(left, right float64, op token.Token) (float64, error) {
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

func resolveStringMath(left string, right any, op token.Token) (any, error) {
	switch right := right.(type) {
	case int:
		if iLeft, err := strconv.Atoi(left); err == nil {
			return resolveIntsMath(iLeft, right, op), nil
		}

		if fLeft, err := strconv.ParseFloat(left, 64); err == nil {
			return resolveFloatsMath(fLeft, float64(right), op)
		}

		if op == token.MUL {
			return strings.Repeat(left, right), nil
		}

		return resolveStringsMath(left, strconv.Itoa(right), op)
	case float64:
		if fLeft, err := strconv.ParseFloat(left, 64); err == nil {
			return resolveFloatsMath(fLeft, right, op)
		}

		return resolveStringsMath(left, strconv.FormatFloat(right, 'f', -1, 64), op)
	case string:
		return resolveStringsMath(left, right, op)
	case []byte:
		return resolveBinariesMath([]byte(left), right, op)
	case []any:
		return resolveSlicesMath([]any{left}, right, op)
	default:
		return nil, fmt.Errorf("right (%T): %w", right, errors.ErrInvalidTypeForMathematicalOperation)
	}
}

func resolveStringsMath(left, right string, op token.Token) (string, error) {
	switch op {
	case token.ADD:
		return left + right, nil
	default:
		return "", fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithStrings)
	}
}

func resolveBinaryMath(left []byte, right any, op token.Token) (any, error) {
	switch right := right.(type) {
	case int:
		if iLeft, err := strconv.Atoi(string(left)); err == nil {
			return resolveIntsMath(iLeft, right, op), nil
		}

		if fLeft, err := strconv.ParseFloat(string(left), 64); err == nil {
			return resolveFloatsMath(fLeft, float64(right), op)
		}

		if op == token.MUL {
			return bytes.Repeat(left, right), nil
		}

		return resolveBinariesMath(left, []byte(strconv.Itoa(right)), op)
	case float64:
		if fLeft, err := strconv.ParseFloat(string(left), 64); err == nil {
			return resolveFloatsMath(fLeft, right, op)
		}

		return resolveBinariesMath(left, []byte(strconv.FormatFloat(right, 'f', -1, 64)), op)
	case string:
		return resolveBinariesMath(left, []byte(right), op)
	case []byte:
		return resolveBinariesMath(left, right, op)
	case []any:
		return resolveSlicesMath([]any{left}, right, op)

	default:
		return nil, fmt.Errorf("right (%T): %w", right, errors.ErrInvalidTypeForMathematicalOperation)
	}
}

func resolveBinariesMath(left, right []byte, op token.Token) ([]byte, error) {
	switch op {
	case token.ADD:
		return append(left, right...), nil
	default:
		return nil, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithBinaries)
	}
}

func resolveSliceMath(left []any, right any, op token.Token) (any, error) {
	switch right := right.(type) {
	case int:
		if op == token.MUL {
			return slices.Repeat(left, right), nil
		}

		return resolveSlicesMath(left, []any{right}, op)
	case float64, string, []byte:
		return resolveSlicesMath(left, []any{right}, op)
	case []any:
		return resolveSlicesMath(left, right, op)
	default:
		return nil, fmt.Errorf("right (%T): %w", right, errors.ErrInvalidTypeForMathematicalOperation)
	}
}

func resolveSlicesMath(left, right []any, op token.Token) ([]any, error) {
	if op == token.ADD {
		result := make([]any, 0, len(left)+len(right))
		result = append(result, left...)
		result = append(result, right...)

		return result, nil
	}

	return nil, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithLists)
}

func resolveMapMath(left map[string]any, right any, op token.Token) (map[string]any, error) {
	mRight, ok := right.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("right (%T): %w", right, errors.ErrInvalidTypeForMathematicalOperation)
	}

	if op == token.ADD {
		result := make(map[string]any, len(left)+len(mRight))
		maps.Copy(result, left)
		maps.Copy(result, mRight)

		return result, nil
	}

	return nil, fmt.Errorf("%q: %w", op.String(), errors.ErrInvalidOperationWithMaps)
}
