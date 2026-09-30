package convert

// AnyToBool converts any value to a bool.
func AnyToBool(value any) bool {
	switch value := value.(type) {
	case bool:
		return value
	case int:
		return value != 0
	case int8:
		return value != 0
	case int16:
		return value != 0
	case int32:
		return value != 0
	case int64:
		return value != 0
	case float32:
		return value != 0
	case float64:
		return value != 0
	case string:
		return value != ""
	case []any:
		return len(value) != 0
	case map[string]any:
		return len(value) != 0
	default:
		return false
	}
}
