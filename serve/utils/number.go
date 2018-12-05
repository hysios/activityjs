package utils

import "strconv"

// IsFloat 宽松的转换大多数类型到整型类型
func IsInt(val interface{}) (int, bool) {
	switch v := val.(type) {
	case int32:
		return int(v), true

	case int64:
		return int(v), true
	case int:
		return v, true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	case bool:
		if v {
			return 1, true
		} else {
			return 0, true
		}
	case string:
		if n, err := strconv.Atoi(v); err != nil {
			return 0, false
		} else {
			return n, true
		}
	default:
		return 0, false
	}
}

// IsFloat 宽松的转换大多数类型到浮点类型
func IsFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return float64(v), true
	case bool:
		if v {
			return 1.0, true
		} else {
			return 0, true
		}
	case string:
		if n, err := strconv.ParseFloat(v, 64); err != nil {
			return 0, false
		} else {
			return n, true
		}
	default:
		return 0, false
	}
}
