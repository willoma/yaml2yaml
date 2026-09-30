package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// AnyToString converts any value to a string.
func AnyToString(value any) string {
	switch value := value.(type) {
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case string:
		return value
	case []any:
		result := make([]string, len(value))
		for i, item := range value {
			result[i] = AnyToString(item)
		}

		return strings.Join(result, ", ")
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", value)
	}
}
