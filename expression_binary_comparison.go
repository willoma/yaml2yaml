package yaml2yaml

import (
	"bytes"
	"go/token"
	"strconv"
)

func resolveComparison(left, right any, op token.Token) bool {
	switch left := left.(type) {
	case bool:
		return resolveBoolComparison(left, right, op)
	case int:
		return resolveIntComparison(left, right, op)
	case float64:
		return resolveFloatComparison(left, right, op)
	case string:
		return resolveStringComparison(left, right, op)
	case []byte:
		return resolveBinaryComparison(left, right, op)
	case []any:
		return resolveSliceComparison(left, right, op)
	case map[string]any:
		return resolveMapComparison(left, right, op)
	default:
		return false
	}
}

func resolveBoolComparison(left bool, right any, op token.Token) bool {
	switch right := right.(type) {
	case bool:
		return resolveBoolsComparison(left, right, op)
	case int:
		return resolveBoolsComparison(left, right > 0, op)
	case float64:
		return resolveBoolsComparison(left, right > 0, op)
	case string:
		if bRight, err := strconv.ParseBool(right); err == nil {
			return resolveBoolsComparison(left, bRight, op)
		}

		return resolveTypedComparison(strconv.FormatBool(left), right, op)
	default:
		return op == token.NEQ
	}
}

func resolveBoolsComparison(left, right bool, op token.Token) bool {
	switch op {
	case token.EQL, token.LEQ, token.GEQ:
		return left == right
	case token.NEQ:
		return left != right
	default:
		return false
	}
}
func resolveIntComparison(left int, right any, op token.Token) bool {
	switch right := right.(type) {
	case bool:
		return resolveBoolsComparison(left > 0, right, op)
	case int:
		return resolveTypedComparison(left, right, op)
	case float64:
		return resolveTypedComparison(float64(left), right, op)
	case string:
		if iRight, err := strconv.Atoi(right); err == nil {
			return resolveTypedComparison(left, iRight, op)
		}

		if fRight, err := strconv.ParseFloat(right, 64); err == nil {
			return resolveTypedComparison(float64(left), fRight, op)
		}

		return resolveTypedComparison(strconv.Itoa(left), right, op)
	default:
		return op == token.NEQ
	}
}

func resolveFloatComparison(left float64, right any, op token.Token) bool {
	switch right := right.(type) {
	case bool:
		return resolveBoolsComparison(left > 0, right, op)
	case int:
		return resolveTypedComparison(left, float64(right), op)
	case float64:
		return resolveTypedComparison(left, right, op)
	case string:
		if fRight, err := strconv.ParseFloat(right, 64); err == nil {
			return resolveTypedComparison(left, fRight, op)
		}

		return resolveTypedComparison(strconv.FormatFloat(left, 'f', -1, 64), right, op)
	default:
		return op == token.NEQ
	}
}

func resolveStringComparison(left string, right any, op token.Token) bool {
	switch right := right.(type) {
	case bool:
		if bLeft, err := strconv.ParseBool(left); err == nil {
			return resolveBoolsComparison(bLeft, right, op)
		}

		return resolveTypedComparison(left, strconv.FormatBool(right), op)
	case int:
		if iLeft, err := strconv.Atoi(left); err == nil {
			return resolveTypedComparison(iLeft, right, op)
		}

		if fLeft, err := strconv.ParseFloat(left, 64); err == nil {
			return resolveTypedComparison(fLeft, float64(right), op)
		}

		return resolveTypedComparison(left, strconv.Itoa(right), op)
	case float64:
		if fLeft, err := strconv.ParseFloat(left, 64); err == nil {
			return resolveTypedComparison(fLeft, right, op)
		}

		return resolveTypedComparison(left, strconv.FormatFloat(right, 'f', -1, 64), op)
	case string:
		return resolveTypedComparison(left, right, op)
	default:
		return op == token.NEQ
	}
}

func resolveTypedComparison[T ~int | ~float64 | ~string](left, right T, op token.Token) bool {
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

func resolveBinaryComparison(left []byte, right any, op token.Token) bool {
	switch right := right.(type) {
	case bool:
		if bLeft, err := strconv.ParseBool(string(left)); err == nil {
			return resolveBoolsComparison(bLeft, right, op)
		}

		return op == token.NEQ
	case string:
		return resolveTypedComparison(string(left), right, op)
	case []byte:
		result := bytes.Compare(left, right)

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
	default:
		return op == token.NEQ
	}
}

func resolveSliceComparison(left []any, right any, op token.Token) bool {
	sRight, ok := right.([]any)
	if !ok {
		return op == token.NEQ
	}

	switch {
	case len(left) > len(sRight):
		return op == token.NEQ || op == token.GTR || op == token.GEQ
	case len(left) < len(sRight):
		return op == token.NEQ || op == token.LSS || op == token.LEQ
	}

	for i := range left {
		if !resolveComparison(left[i], sRight[i], op) {
			return false
		}
	}

	return true
}

func resolveMapComparison(left map[string]any, right any, op token.Token) bool {
	mRight, ok := right.(map[string]any)
	if !ok {
		return op == token.NEQ
	}

	switch {
	case len(left) > len(mRight):
		return op == token.NEQ || op == token.GTR || op == token.GEQ
	case len(left) < len(mRight):
		return op == token.NEQ || op == token.LSS || op == token.LEQ
	}

	for k := range left {
		if !resolveComparison(left[k], mRight[k], op) {
			return false
		}
	}

	return true
}
